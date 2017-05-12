package tests

import (
	"github.com/go-openapi/strfmt"
)

//go:generate msgp
//msgp:shim strfmt.DateTime as:int64 using:dateTimeToInt64/int64ToDateTime
//msgp:shim strfmt.Email as:string using:(strfmt.Email).String/strfmt.Email
import (
	"math/rand"
	"sort"
	"sync"
	"time"
)

type PVersion uint32
type PHash uint32

//msgp:ignore PerfData
type PerfData struct {
	mutex sync.RWMutex

	Status  *PStatus
	users   []*PUser
	forums  []*PForum
	threads []*PThread
	posts   []*PPost

	lastIndex int32

	threadsByForum    map[string][]*PThread
	usersByForum      map[string]map[*PUser]bool
	postsByThreadFlat map[int32][]*PPost
	postsByThreadTree map[int32][]*PPost
	userByNickname    map[string]*PUser
	forumBySlug       map[string]*PForum
	postById          map[int64]*PPost
	threadById        map[int32]*PThread
	threadBySlug      map[string]*PThread
}

//msgp:ignore PStatus
type PStatus struct {
	Version PVersion
	Forum   int32
	Post    int64
	Thread  int32
	User    int32
}

type PUser struct {
	Version      PVersion     `msg:"-"`
	AboutHash    PHash        `msg:"about"`
	Email        strfmt.Email `msg:"email"`
	FullnameHash PHash        `msg:"name"`
	Nickname     string       `msg:"nick"`
}

type PThread struct {
	mutex sync.RWMutex

	Version     PVersion        `msg:"-"`
	ID          int32           `msg:"id"`
	Slug        string          `msg:"slug"`
	Author      *PUser          `msg:"-"`
	Forum       *PForum         `msg:"-"`
	MessageHash PHash           `msg:"message"`
	TitleHash   PHash           `msg:"title"`
	Created     strfmt.DateTime `msg:"created"`
	Votes       int32           `msg:"-"`
	Posts       int32           `msg:"-"`
}

type PForum struct {
	Version   PVersion `msg:"-"`
	Posts     int64    `msg:"-"`
	Slug      string   `msg:"slug"`
	Threads   int32    `msg:"-"`
	TitleHash PHash    `msg:"title"`
	User      *PUser   `msg:"-"`
}

type PPost struct {
	Version     PVersion        `msg:"-"`
	ID          int64           `msg:"id"`
	Author      *PUser          `msg:"-"`
	Thread      *PThread        `msg:"-"`
	Parent      *PPost          `msg:"-"`
	Created     strfmt.DateTime `msg:"created"`
	IsEdited    bool            `msg:"edited"`
	MessageHash PHash           `msg:"message"`
	Index       int32           `msg:"-"`
	Path        []int32         `msg:"-"`
}

func NewPerfData(config *PerfConfig) *PerfData {
	return &PerfData{
		Status:            &PStatus{},
		forums:            make([]*PForum, 0, config.ForumCount),
		users:             make([]*PUser, 0, config.UserCount),
		threads:           make([]*PThread, 0, config.ThreadCount),
		posts:             make([]*PPost, 0, config.PostCount),
		threadsByForum:    map[string][]*PThread{},
		usersByForum:      map[string]map[*PUser]bool{},
		postsByThreadTree: map[int32][]*PPost{},
		postsByThreadFlat: map[int32][]*PPost{},
		userByNickname:    map[string]*PUser{},
		forumBySlug:       map[string]*PForum{},
		threadBySlug:      map[string]*PThread{},
		threadById:        map[int32]*PThread{},
		postById:          map[int64]*PPost{},
	}
}

func (self *PerfData) GetForumBySlug(slug string) *PForum {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.forumBySlug[slug]
}

func (self *PerfData) GetUserByNickname(nickname string) *PUser {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.userByNickname[nickname]
}

func getRandomIndex(count int) int {
	if count == 0 {
		return -1
	}
	return rand.Intn(count)
}

func (self *PerfData) AddForum(forum *PForum) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if _, ok := self.forumBySlug[forum.Slug]; ok {
		panic("Internal error: forum.Slug = " + forum.Slug)
	}
	self.forums = append(self.forums, forum)
	self.forumBySlug[forum.Slug] = forum
	self.usersByForum[forum.Slug] = map[*PUser]bool{}
	self.threadsByForum[forum.Slug] = []*PThread{}
	self.Status.Forum++
}

func (self *PerfData) GetForum(index int) *PForum {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.forums))
	}
	return self.forums[index]
}

func (self *PerfData) AddUser(user *PUser) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if _, ok := self.userByNickname[user.Nickname]; ok {
		panic("Internal error: user.Nickname = " + user.Nickname)
	}
	self.users = append(self.users, user)
	self.userByNickname[user.Nickname] = user
	self.Status.User++
}

func (self *PerfData) GetUser(index int) *PUser {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.users))
	}
	return self.users[index]
}

func (self *PerfData) AddThread(thread *PThread) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if thread.Slug != "" {
		if _, ok := self.threadBySlug[thread.Slug]; ok {
			panic("Internal error: thread.Slug = " + thread.Slug)
		}
		self.threadBySlug[thread.Slug] = thread
	}
	self.threads = append(self.threads, thread)
	self.threadById[thread.ID] = thread
	self.postsByThreadTree[thread.ID] = []*PPost{}
	self.postsByThreadFlat[thread.ID] = []*PPost{}
	self.threadsByForum[thread.Forum.Slug] = append(self.threadsByForum[thread.Forum.Slug], thread)
	self.usersByForum[thread.Forum.Slug][thread.Author] = true
	thread.Forum.Threads++
	self.Status.Thread++
}

func (self *PerfData) GetThread(index int) *PThread {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.threads))
	}
	return self.threads[index]
}

func (self *PerfData) GetThreadById(id int32) *PThread {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.threadById[id]
}

func (self *PerfData) GetPostById(id int64) *PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postById[id]
}

func (self *PerfData) GetForumThreads(forum *PForum) []*PThread {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	result := []*PThread{}
	if result != nil {
		result = append(result, self.threadsByForum[forum.Slug]...)
	}
	return result
}

func (self *PerfData) GetForumUsers(forum *PForum) []*PUser {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	users := self.usersByForum[forum.Slug]
	if users == nil {
		return []*PUser{}
	}
	result := make([]*PUser, 0, len(users))
	for k := range users {
		result = append(result, k)
	}
	return result
}

func (self *PerfData) GetThreadPostsFlat(thread *PThread) []*PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postsByThreadFlat[thread.ID]
}

func (self *PerfData) GetThreadPostsTree(thread *PThread) []*PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postsByThreadTree[thread.ID]
}

func (self *PerfData) AddPost(post *PPost) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.posts = append(self.posts, post)
	self.postById[post.ID] = post
	self.usersByForum[post.Thread.Forum.Slug][post.Author] = true

	self.lastIndex++
	post.Index = self.lastIndex
	if post.Parent != nil {
		// Явное копирование массива, т.к. append не всегда ведёт себя адекватно в многопоточном окружении
		path := make([]int32, 0, len(post.Parent.Path)+1)
		path = append(path, post.Parent.Path...)
		path = append(path, post.Index)
		post.Path = path
	} else {
		post.Path = []int32{post.Index}
	}

	tree := append(self.postsByThreadTree[post.Thread.ID], post)
	self.postsByThreadFlat[post.Thread.ID] = append(self.postsByThreadFlat[post.Thread.ID], post)
	if post.Parent != nil {
		sort.Sort(PPostSortTree(tree))
	}
	self.postsByThreadTree[post.Thread.ID] = tree

	post.Thread.Forum.Posts++
	post.Thread.Posts++
	self.Status.Post++
}

func (self *PerfData) GetPost(index int) *PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.posts))
	}
	return self.posts[index]
}

func (self *PPost) GetParentId() int64 {
	if self.Parent == nil {
		return 0
	}
	return self.Parent.ID
}

func GetRandomLimit() int32 {
	return 15 + rand.Int31n(5)
}

func GetRandomDesc() *bool {
	switch rand.Intn(3) {
	case 0:
		v := false
		return &v
	case 1:
		v := true
		return &v
	default:
		return nil
	}
}

func dateTimeToInt64(value strfmt.DateTime) int64 {
	return time.Time(value).UnixNano()
}

func int64ToDateTime(value int64) strfmt.DateTime {
	return strfmt.DateTime(time.Unix(0, value))
}