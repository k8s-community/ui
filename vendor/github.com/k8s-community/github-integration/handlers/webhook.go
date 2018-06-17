package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/google/go-github/github"
	"github.com/k8s-community/cicd"
	"github.com/k8s-community/github-integration/models"
	userManClient "github.com/k8s-community/user-manager/client"
	"github.com/takama/router"
	"gopkg.in/reform.v1"
	githubhook "gopkg.in/rjz/githubhook.v0"
)

// WebHookHandler is common handler for web hooks (installation, repositories installation, push)
func (h *Handler) WebHookHandler(c *router.Control) {
	secret := []byte(h.Env["GITHUBINT_TOKEN"])

	hook, err := githubhook.Parse(secret, c.Request)
	if err != nil {
		h.Errlog.Printf("cannot parse hook (ID %s): %s", hook.Id, err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		// Triggered when an integration has been installed or uninstalled by user.
		h.Infolog.Printf("initialization web hook (ID %s)", hook.Id)
		err = h.saveInstallation(hook)

	case "integration_installation_repositories":
		// Triggered when a repository is added or removed from an installation.
		h.Infolog.Printf("initialization web hook for user repositories (ID %s)", hook.Id)
		err = h.initialUserManagement(hook)

	case "push":
		// Any Git push to a Repository, including editing tags or branches.
		// Commits via API actions that update references are also counted. This is the default event.
		h.Infolog.Printf("push hook (ID %s)", hook.Id)
		err = h.processPush(c, hook)
		if err != nil {
			h.Infolog.Printf("cannot run ci/cd process for hook (ID %s): %s", hook.Id, err)
			c.Code(http.StatusBadRequest).Body(nil)
			return
		}

	case "create":
		h.Infolog.Printf("create hook (ID %s)", hook.Id)
		// ToDo: keep it for the future
		/*err = h.processCreate(c, hook)
		if err != nil {
			h.Infolog.Printf("cannot run ci/cd process for hook (ID %s): %s", hook.Id, err)
			c.Code(http.StatusBadRequest).Body(nil)
			return
		}*/
		return

	default:
		h.Infolog.Printf("Warning! Don't know how to process hook (ID %s), event = %s", hook.Id, hook.Event)
		c.Code(http.StatusOK).Body(nil)
		return
	}

	if err != nil {
		h.Errlog.Printf("cannot process hook (ID %s, event = %s): %s", hook.Id, hook.Event, err)
		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	h.Infolog.Printf("finished to process hook (ID %s, event = %s)", hook.Id, hook.Event)
	c.Code(http.StatusOK).Body(nil)
}

// initialUserManagement is used for user activation in k8s system
func (h *Handler) initialUserManagement(hook *githubhook.Hook) error {
	evt := github.InstallationRepositoriesEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	userManagerURL := h.Env["USERMAN_BASE_URL"]

	client, err := userManClient.NewClient(nil, userManagerURL)
	if err != nil {
		return err
	}

	h.Infolog.Print("Try to activate (sync) user in k8s system: ", *evt.Sender.Login)

	user := userManClient.NewUser(*evt.Installation.Account.Login)

	code, err := client.User.Sync(user)
	if err != nil {
		return err
	}

	h.Infolog.Printf("Service user-man, method sync, returned code: %d", code)

	return nil
}

// processPush is used for start CI/CD process for some repository from push hook
func (h *Handler) processPush(c *router.Control, hook *githubhook.Hook) error {
	evt := github.PushEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	// ToDO: process somehow kind of hooks without HeadCommit
	if evt.HeadCommit == nil {
		h.Infolog.Printf("Warning! Don't know how to process hook %s - no HeadCommit inside", hook.Id)
		return nil
	}

	h.setInstallationID(*evt.Repo.Owner.Name, *evt.Installation.ID)

	prefix := "refs/heads/" + h.Env["GITHUBINT_BRANCH"]
	if !strings.HasPrefix(*evt.Ref, prefix) {
		h.Infolog.Printf("Warning! Don't know how to process hook %s - branch %s", hook.Id, *evt.Ref)
		return nil
	}

	ciCdURL := h.Env["CICD_BASE_URL"]

	client := cicd.NewClient(ciCdURL)

	version := strings.Trim(*evt.Ref, prefix)

	// run CICD process
	req := &cicd.BuildRequest{
		Username:   *evt.Repo.Owner.Name,
		Repository: *evt.Repo.Name,
		CommitHash: *evt.HeadCommit.ID,
		Task:       cicd.TaskDeploy,
		Version:    &version,
	}

	_, err = client.Build(req)
	if err != nil {
		return fmt.Errorf("cannot run ci/cd process for hook (ID %s): %s", hook.Id, err)
	}

	return nil
}

func (h *Handler) processCreate(c *router.Control, hook *githubhook.Hook) error {
	evt := github.CreateEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	// Process only tags
	if evt.RefType == nil || *evt.RefType != "tag" {
		h.Infolog.Printf("Warning! Don't know how to process hook %s - not a tag", hook.Id)
		return nil
	}

	h.setInstallationID(*evt.Repo.Owner.Name, *evt.Installation.ID)

	ciCdURL := h.Env["CICD_BASE_URL"]

	client := cicd.NewClient(ciCdURL)

	// run CICD process
	req := &cicd.BuildRequest{
		Username:   *evt.Repo.Owner.Name,
		Repository: *evt.Repo.Name,
		CommitHash: *evt.Ref,
		Task:       cicd.TaskDeploy,
		Version:    evt.Ref,
	}

	_, err = client.Build(req)
	if err != nil {
		return fmt.Errorf("cannot run ci/cd process for hook (ID %s): %s", hook.Id, err)
	}

	return nil
}

// saveInstallation saves installation in memory
func (h *Handler) saveInstallation(hook *githubhook.Hook) error {
	evt := github.InstallationEvent{}

	err := hook.Extract(&evt)
	if err != nil {
		return err
	}

	h.Infolog.Printf("save installation for user %s (installation ID = %d)", *evt.Sender.Login, *evt.Installation.ID)

	// save installation for commit status update
	err = h.setInstallationID(*evt.Installation.Account.Login, *evt.Installation.ID)
	if err != nil {
		h.Errlog.Printf("Couldn't save installation: %+v", err)
	}

	return nil
}

// installationID gets installation from DB
func (h *Handler) installationID(username string) (*int, error) {
	st, err := h.DB.FindOneFrom(models.InstallationTable, "username", username)
	if err != nil {
		return nil, err
	}
	inst := st.(*models.Installation)

	return pointer.ToInt(inst.InstallationID), nil
}

func (h *Handler) setInstallationID(username string, instID int) error {
	var inst = &models.Installation{}

	st, err := h.DB.FindOneFrom(models.InstallationTable, "username", username)
	if err != nil && err != reform.ErrNoRows {
		return err
	}

	if err == nil {
		inst = st.(*models.Installation)
	}

	inst.InstallationID = instID
	inst.Username = username

	err = h.DB.Save(inst)

	return err
}
