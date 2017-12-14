// Copyright 2016-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package inmemory

//
// import (
// 	"fmt"
// 	"sync"
//
// 	"github.com/omarghader/jobqueue/dao"
// 	"github.com/omarghader/jobqueue/models"
// 	"github.com/omarghader/jobqueue/protocol"
// )
//
// // InMemoryStore is a simple in-memory store implementation.
// // It implements the Store interface. Do not use in production.
// type InMemoryStore struct {
// 	mu   sync.Mutex
// 	jobs map[string]*models.Job
// }
//
// // NewInMemoryStore creates a new InMemoryStore.
// func NewInMemoryStore() *InMemoryStore {
// 	return &InMemoryStore{
// 		jobs: make(map[string]*models.Job),
// 	}
// }
//
// // Start the store.
// func (st *InMemoryStore) Start() error {
// 	return nil
// }
//
// // Create adds a new job.
// func (st *InMemoryStore) Create(job *protocol.Job) error {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	st.jobs[job.ID] = job
// 	return nil
// }
//
// // Delete removes the job.
// func (st *InMemoryStore) Delete(job *protocol.Job) error {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	delete(st.jobs, job.ID)
// 	return nil
// }
//
// // Update updates the job.
// func (st *InMemoryStore) Update(job *protocol.Job) error {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	st.jobs[job.ID] = job
// 	return nil
// }
//
// // Next picks the next job to execute.
// func (st *InMemoryStore) Next() (*protocol.Job, error) {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	var next *models.Job
// 	for _, job := range st.jobs {
// 		if job.State == models.Waiting {
// 			if next == nil || job.Rank > next.Rank || job.Priority > next.Priority {
// 				next = job
// 			}
// 		}
// 	}
// 	return next, nil
// }
//
// // Stats returns statistics about the jobs in the store.
// func (st *InMemoryStore) Stats(req *models.StatsRequest) (*models.Stats, error) {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	stats := &models.Stats{}
// 	for _, job := range st.jobs {
// 		if req.Topic != "" && job.Topic != req.Topic {
// 			continue
// 		}
// 		if req.CorrelationGroup != "" && job.CorrelationGroup != req.CorrelationGroup {
// 			continue
// 		}
// 		switch job.State {
// 		default:
// 			return nil, fmt.Errorf("found unknown state %v", job.State)
// 		case models.Waiting:
// 			stats.Waiting++
// 		case models.Working:
// 			stats.Working++
// 		case models.Succeeded:
// 			stats.Succeeded++
// 		case models.Failed:
// 			stats.Failed++
// 		}
// 	}
// 	return stats, nil
// }
//
// // Lookup returns the job with the specified identifier (or ErrNotFound).
// func (st *InMemoryStore) Lookup(id string) (*models.Job, error) {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	job, found := st.jobs[id]
// 	if !found {
// 		return nil, dao.ErrNotFound
// 	}
// 	return job, nil
// }
//
// // LookupByCorrelationID returns the details of jobs by their correlation identifier.
// // If no such job could be found, an empty array is returned.
// func (st *InMemoryStore) LookupByCorrelationID(correlationID string) ([]*models.Job, error) {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	var result []*models.Job
// 	for _, job := range st.jobs {
// 		if job.CorrelationID == correlationID {
// 			result = append(result, job)
// 		}
// 	}
// 	return result, nil
// }
//
// // List finds matching jobs.
// func (st *InMemoryStore) List(req *protocol.ListRequest) (*protocol.ListResponse, error) {
// 	st.mu.Lock()
// 	defer st.mu.Unlock()
// 	rsp := &protocol.ListResponse{}
// 	i := -1
// 	for _, job := range st.jobs {
// 		var skip bool
// 		if req.Topic != "" && req.Topic != job.Topic {
// 			skip = true
// 		}
// 		if req.State != "" && job.State != req.State {
// 			skip = true
// 		}
// 		if req.Offset > 0 && req.Offset < i {
// 			skip = true
// 			rsp.Total++
// 		}
// 		if req.Limit > 0 && i > req.Limit {
// 			skip = true
// 			rsp.Total++
// 		}
// 		if !skip {
// 			j, _ := job.ToJob()
// 			rsp.Jobs = append(rsp.Jobs, *j)
// 		}
// 		i++
// 	}
// 	return rsp, nil
// }
