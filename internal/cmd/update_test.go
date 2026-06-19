package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestUpdate(t *testing.T) {
	suite.Run(t, &UpdateTestSuite{})
}

type UpdateTestSuite struct {
	suite.Suite
}

func (suite *UpdateTestSuite) TestUpdateAlgorithm() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN", Algorithm: "sha1"}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{Algorithm: "sha256"}, []string{"algorithm"}, false, false)

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	acc, _ := ns.FindAccount("github")
	suite.Equal("sha256", acc.Algorithm)
}

func (suite *UpdateTestSuite) TestUpdateLength() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN", Length: 6}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{Length: 8}, []string{"length"}, false, false)

	suite.Require().NoError(err)

	acc, _ := ns.FindAccount("github")
	suite.Equal(uint(8), acc.Length)
}

func (suite *UpdateTestSuite) TestUpdatePrefix() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN"}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{Prefix: "GH:"}, []string{"prefix"}, false, false)

	suite.Require().NoError(err)

	acc, _ := ns.FindAccount("github")
	suite.Equal("GH:", acc.Prefix)
}

func (suite *UpdateTestSuite) TestUpdateTimePeriod() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN", TimePeriod: 30}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{TimePeriod: 60}, []string{"time-period"}, false, false)

	suite.Require().NoError(err)

	acc, _ := ns.FindAccount("github")
	suite.Equal(int64(60), acc.TimePeriod)
}

func (suite *UpdateTestSuite) TestUpdateOnlyFlaggedFields() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN", Algorithm: "sha1", Length: 6}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{Algorithm: "sha256", Length: 8}, []string{"algorithm"}, false, false)

	suite.Require().NoError(err)

	acc, _ := ns.FindAccount("github")
	suite.Equal("sha256", acc.Algorithm)
	suite.Equal(uint(6), acc.Length)
}

func (suite *UpdateTestSuite) TestUpdateNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeUpdate(stor, "missing", "github", AccountOptions{}, []string{}, false, false)

	suite.Require().Error(err)
}

func (suite *UpdateTestSuite) TestUpdateAccountNotFound() {
	ns := &storage.Namespace{Name: "myns"}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "missing", AccountOptions{}, []string{}, false, false)

	suite.Require().Error(err)
}

func (suite *UpdateTestSuite) TestUpdateSaveError() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN"}},
	}
	stor := newFakeStorage(ns)
	stor.SaveErr = errors.New("disk full")

	err := executeUpdate(stor, "myns", "github", AccountOptions{}, []string{}, false, false)

	suite.Require().Error(err)
}

func (suite *UpdateTestSuite) TestUpdateSetFavorite() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN"}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{}, []string{}, true, false)

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	acc, _ := ns.FindAccount("github")
	suite.True(acc.Favorite)
}

func (suite *UpdateTestSuite) TestUpdateUnsetFavorite() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "TOKEN", Favorite: true}},
	}
	stor := newFakeStorage(ns)

	err := executeUpdate(stor, "myns", "github", AccountOptions{}, []string{}, false, true)

	suite.Require().NoError(err)

	acc, _ := ns.FindAccount("github")
	suite.False(acc.Favorite)
}
