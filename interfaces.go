package raven

// Message ... https://docs.getsentry.com/hosted/clientdev/interfaces/#message-interface
type Message struct {
	// Required
	Message string `json:"message"`

	// Optional
	Params []interface{} `json:"params,omitempty"`
}

// Class ...
func (m *Message) Class() string { return "logentry" }

// Template ... https://docs.getsentry.com/hosted/clientdev/interfaces/#template-interface
type Template struct {
	// Required
	Filename    string `json:"filename"`
	Lineno      int    `json:"lineno"`
	ContextLine string `json:"context_line"`

	// Optional
	PreContext   []string `json:"pre_context,omitempty"`
	PostContext  []string `json:"post_context,omitempty"`
	AbsolutePath string   `json:"abs_path,omitempty"`
}

// Class ...
func (t *Template) Class() string { return "template" }

// User ... https://docs.getsentry.com/hosted/clientdev/interfaces/#context-interfaces
type User struct {
	// All fields are optional
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	IP       string `json:"ip_address,omitempty"`
}

// Class ...
func (h *User) Class() string { return "user" }

// Query ... https://docs.getsentry.com/hosted/clientdev/interfaces/#context-interfaces
type Query struct {
	// Required
	Query string `json:"query"`

	// Optional
	Engine string `json:"engine,omitempty"`
}

// Class ...
func (q *Query) Class() string { return "query" }
