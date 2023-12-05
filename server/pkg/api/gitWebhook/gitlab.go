package gitWebhook

import "time"

// gitlab的Push返回的所有Json
type GitlabPushHook struct {
	ObjectKind   string      `json:"object_kind"`
	EventName    string      `json:"event_name"`
	Before       string      `json:"before"`
	After        string      `json:"after"`
	Ref          string      `json:"ref"`
	CheckoutSha  string      `json:"checkout_sha"`
	Message      interface{} `json:"message"`
	UserId       int         `json:"user_id"`
	UserName     string      `json:"user_name"`
	UserUsername string      `json:"user_username"`
	UserEmail    string      `json:"user_email"`
	UserAvatar   string      `json:"user_avatar"`
	ProjectId    int         `json:"project_id"`
	Project      struct {
		Id                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebUrl            string      `json:"web_url"`
		AvatarUrl         interface{} `json:"avatar_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		CiConfigPath      string      `json:"ci_config_path"`
		Homepage          string      `json:"homepage"`
		Url               string      `json:"url"`
		SshUrl            string      `json:"ssh_url"`
		HttpUrl           string      `json:"http_url"`
	} `json:"project"`
	Commits []struct {
		Id        string    `json:"id"`
		Message   string    `json:"message"`
		Title     string    `json:"title"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []interface{} `json:"added"`
		Modified []string      `json:"modified"`
		Removed  []interface{} `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int `json:"total_commits_count"`
	PushOptions       struct {
	} `json:"push_options"`
	Repository struct {
		Name            string `json:"name"`
		Url             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitHttpUrl      string `json:"git_http_url"`
		GitSshUrl       string `json:"git_ssh_url"`
		VisibilityLevel int    `json:"visibility_level"`
	} `json:"repository"`
}

// gitlab的Merge Request Hook返回的所有Json
type GitlabPRHook struct {
	ObjectKind string `json:"object_kind"`
	EventType  string `json:"event_type"`
	User       struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarUrl string `json:"avatar_url"`
		Email     string `json:"email"`
	} `json:"user"`
	Project struct {
		Id                int         `json:"id"`
		Name              string      `json:"name"`
		Description       string      `json:"description"`
		WebUrl            string      `json:"web_url"`
		AvatarUrl         interface{} `json:"avatar_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		Namespace         string      `json:"namespace"`
		VisibilityLevel   int         `json:"visibility_level"`
		PathWithNamespace string      `json:"path_with_namespace"`
		DefaultBranch     string      `json:"default_branch"`
		CiConfigPath      string      `json:"ci_config_path"`
		Homepage          string      `json:"homepage"`
		Url               string      `json:"url"`
		SshUrl            string      `json:"ssh_url"`
		HttpUrl           string      `json:"http_url"`
	} `json:"project"`
	ObjectAttributes struct {
		AssigneeId     interface{} `json:"assignee_id"`
		AuthorId       int         `json:"author_id"`
		CreatedAt      string      `json:"created_at"`
		Description    string      `json:"description"`
		HeadPipelineId interface{} `json:"head_pipeline_id"`
		Id             int         `json:"id"`
		Iid            int64       `json:"iid"`
		LastEditedAt   interface{} `json:"last_edited_at"`
		LastEditedById interface{} `json:"last_edited_by_id"`
		MergeCommitSha interface{} `json:"merge_commit_sha"`
		MergeError     interface{} `json:"merge_error"`
		MergeParams    struct {
			ForceRemoveSourceBranch string `json:"force_remove_source_branch"`
		} `json:"merge_params"`
		MergeStatus               string      `json:"merge_status"`
		MergeUserId               interface{} `json:"merge_user_id"`
		MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
		MilestoneId               interface{} `json:"milestone_id"`
		SourceBranch              string      `json:"source_branch"`
		SourceProjectId           int         `json:"source_project_id"`
		StateId                   int         `json:"state_id"`
		TargetBranch              string      `json:"target_branch"`
		TargetProjectId           int         `json:"target_project_id"`
		TimeEstimate              int         `json:"time_estimate"`
		Title                     string      `json:"title"`
		UpdatedAt                 string      `json:"updated_at"`
		UpdatedById               interface{} `json:"updated_by_id"`
		Url                       string      `json:"url"`
		Source                    struct {
			Id                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebUrl            string      `json:"web_url"`
			AvatarUrl         interface{} `json:"avatar_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      string      `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"source"`
		Target struct {
			Id                int         `json:"id"`
			Name              string      `json:"name"`
			Description       string      `json:"description"`
			WebUrl            string      `json:"web_url"`
			AvatarUrl         interface{} `json:"avatar_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			Namespace         string      `json:"namespace"`
			VisibilityLevel   int         `json:"visibility_level"`
			PathWithNamespace string      `json:"path_with_namespace"`
			DefaultBranch     string      `json:"default_branch"`
			CiConfigPath      string      `json:"ci_config_path"`
			Homepage          string      `json:"homepage"`
			Url               string      `json:"url"`
			SshUrl            string      `json:"ssh_url"`
			HttpUrl           string      `json:"http_url"`
		} `json:"target"`
		LastCommit struct {
			Id        string    `json:"id"`
			Message   string    `json:"message"`
			Title     string    `json:"title"`
			Timestamp time.Time `json:"timestamp"`
			Url       string    `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		WorkInProgress      bool          `json:"work_in_progress"`
		TotalTimeSpent      int           `json:"total_time_spent"`
		TimeChange          int           `json:"time_change"`
		HumanTotalTimeSpent interface{}   `json:"human_total_time_spent"`
		HumanTimeChange     interface{}   `json:"human_time_change"`
		HumanTimeEstimate   interface{}   `json:"human_time_estimate"`
		AssigneeIds         []interface{} `json:"assignee_ids"`
		State               string        `json:"state"`
		Action              string        `json:"action"`
	} `json:"object_attributes"`
	Labels  []interface{} `json:"labels"`
	Changes struct {
		MergeStatus struct {
			Previous string `json:"previous"`
			Current  string `json:"current"`
		} `json:"merge_status"`
	} `json:"changes"`
	Repository struct {
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
}
