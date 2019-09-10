package jiraui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/go-jira/jira"
	"github.com/mitchellh/go-wordwrap"
	"gopkg.in/coryb/yaml.v2"
	ui "gopkg.in/gizak/termui.v2"
)

func countLabelsFromQuery(query string) map[string]int {
	data, _ := runJiraQuery(query)
	return countLabelsFromQueryData(data)
}

func countLabelsFromQueryData(data interface{}) map[string]int {
	counts := make(map[string]int)
	issues := data.(map[string]interface{})["issues"].([]interface{})
	for _, issue := range issues {
		issueLabels := issue.(map[string]interface{})["fields"].(map[string]interface{})["labels"]
		labels := issueLabels.([]interface{})
		if len(labels) == 0 {
			// "NOT LABELLED" isn't a valid label, so no possible conflict here.
			counts["NOT LABELLED"] = counts["NOT LABELLED"] + 1
		} else {
			for _, v := range labels {
				label := v.(string)
				counts[label] = counts[label] + 1
			}
		}
	}
	return counts
}

func RunExternalCommand(fn func() error) error {
	log.Debugf("ShellOut() called with %q", fn)
	deregisterEventHandlers()
	ui.Clear()
	stty := exec.Command("stty", "-f", "/dev/tty", "echo", "opost")
	_ = stty.Run()
	err := fn() // magic happens
	stty = exec.Command("stty", "-f", "/dev/tty", "-echo", "-opost")
	_ = stty.Run()
	registerEventHandlers()
	if err != nil {
		return err
	}
	return nil
}

func runShell() {
	_ = RunExternalCommand(
		func() error {
			cmd := exec.Command("bash")
			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			return cmd.Run()
		})
	changePage()
}

/*
func runJiraCmdEdit(ticketId string) {
	_ = RunExternalCommand(
		func() error {
			opts := getJiraOpts()
			c := jira.NewJira(opts["endpoint"].(string))
			return c.CmdEdit(ticketId)
		})
	switch c := currentPage.(type) {
	case Refresher:
		c.Refresh()
	}
	changePage()
}
*/

/*
func runJiraCmdCreate(project string, summary string) {
	_ = RunExternalCommand(
		func() error {
			opts := getJiraOpts()
			opts["project"] = project
			opts["summary"] = summary
			c := jira.NewJira(opts["endpoint"].(string))
			return c.CmdCreate()
		})
	switch c := currentPage.(type) {
	case Refresher:
		c.Refresh()
	}
	changePage()
}
*/

/*
func runJiraCmdCommentNoEditor(ticketId string, comment string) {
	opts := getJiraOpts()
	opts["comment"] = comment
	c := jira.NewJira(opts["endpoint"].(string))
	c.CmdComment(ticketId)
}

func runJiraCmdAssign(ticketId string, user string) {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	c.CmdAssign(ticketId, user)
}

func runJiraCmdWatch(ticketId string, watcher string, remove bool) {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	if watcher == "" {
		watcher = opts["user"].(string)
	}
	c.CmdWatch(ticketId, watcher, remove)
}

func runJiraCmdVote(ticketId string, up bool) {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	c.CmdVote(ticketId, up)
}

func runJiraCmdLabels(ticketId string, action string, labels []string) {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	err := c.CmdLabels(action, ticketId, labels)
	if err != nil {
		log.Errorf("Error writing labels: %q", err)
	}
}

func runJiraCmdRank(ticketId, targetId string, order jira.RankOrder) {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	err := c.RankIssue(ticketId, targetId, order)
	if err != nil {
		log.Errorf("Error modifying issue rank: %q", err)
	}
}
*/

func findTicketIdInString(line string) string {
	re := regexp.MustCompile(`[A-Z]{2,12}-[0-9]{1,6}`)
	re_too_long := regexp.MustCompile(`[A-Z]{13}-[0-9]{1,6}`)

	if re_too_long.MatchString(line) {
		return ""
	}

	return strings.TrimSpace(re.FindString(line))
}

func runJiraQuery(query string) (interface{}, error) {
	opts := getJiraOpts()
	opts["query"] = query
	c := jira.NewJira(opts["endpoint"].(string))
	//return c.FindIssues()
	log.Infof("TODO: reenable c.FindIssues: %#v", c)
	return nil, nil
}

func JiraQueryAsStrings(query string, templateName string) []string {
	opts := getJiraOpts()
	opts["query"] = query
	c := jira.NewJira(opts["endpoint"].(string))
	//data, _ := c.FindIssues()
	log.Infof("TODO: reenable c.FindIssues: %#v", c)
	var data interface{}
	buf := new(bytes.Buffer)
	if templateName == "" {
		templateName = "jira_ui_list"
	}
	//template := c.GetTemplate(templateName)
	template := ""
	if template == "" {
		template = default_list_template
	}
	//jira.RunTemplate(template, data, buf)
	log.Infof("TODO: reenable c.RunTemplate: %#v, %#v", data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func FetchJiraTicket(id string) (interface{}, error) {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	//return c.ViewIssue(id)
	log.Infof("TODO: reenable c.ViewIssue: %#v", c)
	return nil, nil
}

func JiraTicketAsStrings(data interface{}, templateName string) []string {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	buf := new(bytes.Buffer)
	//template := c.GetTemplate(templateName)
	log.Infof("TODO: reenable c.GetTemplate: %#v", c)
	template := ""
	log.Debugf("JiraTicketsAsStrings: template = %q", template)
	if template == "" {
		template = strings.Replace(default_view_template, "ENDPOINT", opts["endpoint"].(string), 1)
	}
	//jira.RunTemplate(template, data, buf)
	log.Infof("TODO: reenable c.RunTemplate: %#v, %#v", data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func HelpTextAsStrings(data interface{}, templateName string) []string {
	opts := getJiraOpts()
	c := jira.NewJira(opts["endpoint"].(string))
	buf := new(bytes.Buffer)
	//template := c.GetTemplate(templateName)
	log.Infof("TODO: reenable c.GetTemplate: %#v", c)
	template := ""
	if template == "" {
		template = default_help_template
	}
	log.Debugf("HelpTextAsStrings: template = %q", template)
	//jira.RunTemplate(template, data, buf)
	log.Infof("TODO: reenable c.RunTemplate: %#v, %#v", data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func WrapText(lines []string, maxWidth uint) []string {
	out := make([]string, 0)
	insideNoformatBlock := false
	insideCodeBlock := false
	for _, line := range lines {
		if matched, _ := regexp.MatchString(`^\s+\{code`, line); matched {
			insideCodeBlock = !insideCodeBlock
		} else if strings.TrimSpace(line) == "{noformat}" {
			insideNoformatBlock = !insideNoformatBlock
		}
		if maxWidth == 0 || uint(len(line)) < maxWidth || insideCodeBlock || insideNoformatBlock {
			out = append(out, line)
			continue
		}
		if matched, _ := regexp.MatchString(`^[a-z_]+:\s`, line); matched {
			// don't futz with single line field+value.
			// If they are too long, that's their fault.
			out = append(out, line)
			continue
		}
		// wrap text, but preserve indenting
		re := regexp.MustCompile(`^\s*`)
		indenting := re.FindString(line)
		wrappedLines := strings.Split(wordwrap.WrapString(line, maxWidth-uint(len(indenting))), "\n")
		indentedWrappedLines := make([]string, len(wrappedLines))
		for i, wl := range wrappedLines {
			if i == 0 {
				// first line already has the indent
				indentedWrappedLines[i] = wl
			} else {
				indentedWrappedLines[i] = indenting + wl
			}
		}
		out = append(out, indentedWrappedLines...)
	}
	return out
}

func parseYaml(file string, v map[string]interface{}) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debugf("Parsing YAML file: %s", file)
		yaml.Unmarshal(fh, &v)
	}
}

func loadConfigs(opts map[string]interface{}) {
	/*
		paths := figtree.FindParentPaths(".jira.d/jira-ui-config.yml")
		paths = append(figtree.FindParentPaths(".jira.d/config.yml"), paths...)
		paths = append([]string{"/etc/go-jira-ui.yml", "/etc/go-jira.yml"}, paths...)
	*/
	var paths []string

	// iterate paths in reverse
	for i := len(paths) - 1; i >= 0; i-- {
		file := paths[i]
		if _, err := os.Stat(file); err == nil {
			tmp := make(map[string]interface{})
			parseYaml(file, tmp)
			for k, v := range tmp {
				if _, ok := opts[k]; !ok {
					log.Debugf("Setting %q to %#v from %s", k, v, file)
					opts[k] = v
				}
			}
		}
	}
}

func doLogin(opts map[string]interface{}) error {
	c := jira.NewJira(opts["endpoint"].(string))
	fmt.Printf("Logging in as %s:\n", opts["user"])
	//return c.CmdLogin()
	log.Infof("TODO: reenable c.CmdLogin: %#v", c)
	return nil
}

func ensureLoggedIntoJira() error {
	homeDir := os.Getenv("HOME")
	opts := getJiraOpts()
	testSessionQuery := fmt.Sprintf("reporter = %s", opts["user"])
	if _, err := os.Stat(fmt.Sprintf("%s/.jira.d/cookies.js", homeDir)); err != nil {
		return doLogin(opts)
		//} else if data, err := runJiraQuery(testSessionQuery); err != nil {
	} else if _, err := runJiraQuery(testSessionQuery); err != nil {
		return doLogin(opts)
		/*
			} else if val, ok := data.(map[string]interface{})["errorMessages"]; ok {
				if len(val.([]interface{})) > 0 {
					return doLogin(opts)
				}
		*/
	}
	log.Infof("TODO: refactor ensureLoggedIntoJira")
	return nil
}

func getJiraOpts() map[string]interface{} {
	user := os.Getenv("USER")
	home := os.Getenv("HOME")
	defaultQueryFields := "summary,created,updated,priority,status,reporter,assignee,labels"
	defaultSort := "priority asc, created"
	defaultMaxResults := 1000

	opts := make(map[string]interface{})
	defaults := map[string]interface{}{
		"user":        user,
		"endpoint":    os.Getenv("JIRA_ENDPOINT"),
		"queryfields": defaultQueryFields,
		"directory":   fmt.Sprintf("%s/.jira.d/templates", home),
		"sort":        defaultSort,
		"max_results": defaultMaxResults,
		"method":      "GET",
		"quiet":       true,
	}

	for k, v := range cliOpts {
		if _, ok := opts[k]; !ok {
			log.Debugf("Setting %q to %#v from cli options", k, v)
			opts[k] = v
		}
	}

	loadConfigs(opts)
	for k, v := range defaults {
		if _, ok := opts[k]; !ok {
			log.Debugf("Setting %q to %#v from defaults", k, v)
			opts[k] = v
		}
	}
	return opts
}

func lastLineDisplayed(ls *ScrollableList, firstLine int, correction int) int {
	return firstLine + ls.Height - correction
}
