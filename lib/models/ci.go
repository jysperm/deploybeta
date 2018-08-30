package models

import (
	"fmt"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/lib/db"
)

type Commit struct {
	db.ResourceMeta
	Hash     string `json:"hash"`
	AppName  string `json:"appName"`
	Tag      string `json:"tag"`
	StartAt  string `json:"startAt"`
	Duration int64  `json:"duration"`
	Logs     string `json:"logs"`
}

func (commit *Commit) ResourceKey() string {
	return fmt.Sprintf("/apps/%s/versions/%s/commits/%s", commit.App, commit.Version, commit.Hash)
}

func (commit *Commit) Associations() []db.Association {
	return []db.Association{
		commit.App(),
		commit.Version(),
	}
}
func (commit *Commit) App() db.BelongsToAssociation {
	return db.BelongsTo((&Application{Name: commit.AppName}).ResourceKey())
}

func (commit *Commit) Version() db.BelongsToAssociation {
	return db.BelongsTo((&Version{AppName: commit.AppName, Tag: commit.Tag}).ResourceKey())
}
func NewCommitFromVersion(hash string, version *Version) Commit {
	return Commit{
		Hash:    hash,
		AppName: version.AppName,
		Tag:     version.Tag,
	}
}

func (commit *Commit) Create() error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Create(commit)
	})

	if err != nil {
		return errwrap.Wrapf("create commit: {{err}}", err)
	}

	return nil
}

func (commit *Commit) UpdateStartAt(startAt string) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(commit)

		if err != nil {
			tran.SetError(err)
			return
		}

		commit.StartAt = startAt
		tran.Update(commit)
	})

	return err
}

func (commit *Commit) AppendLogs(log string) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(commit)

		if err != nil {
			tran.SetError(err)
			return
		}

		commit.Logs += log
		tran.Update(commit)
	})

	return err
}

func (commit *Commit) UpdateDuration(duration int64) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(commit)

		if err != nil {
			tran.SetError(err)
			return
		}

		commit.Duration = duration
		tran.Update(commit)
	})

	return err
}
