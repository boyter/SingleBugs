package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Flag inputs
var portNumber = flag.Int("port", 8080, "port to bind to (default 8080)")
//var useTemplates = flag.Bool("tmpl", false, "use external templates (default false)")

const dataPath = "data/"
const tmplPath = "tmpl/"
const demoVersion = false
const returnLimit = 50

// Types *NB* names must be uppercase to appear in JSON
type Project struct {
	Id   int64
	Name string
	Open bool
	Time int64
}
type ProjectByTimeDesc []Project

func (this ProjectByTimeDesc) Len() int {
	return len(this)
}
func (this ProjectByTimeDesc) Less(i, j int) bool {
	return this[i].Time > this[j].Time
}
func (this ProjectByTimeDesc) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type Issue struct {
	Id        int64
	Projectid int64
	Name      string
	Open      bool
	Time      int64
}
type IssueByTimeDesc []Issue

func (this IssueByTimeDesc) Len() int {
	return len(this)
}
func (this IssueByTimeDesc) Less(i, j int) bool {
	return this[i].Time > this[j].Time
}
func (this IssueByTimeDesc) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type Note struct {
	Id      int64
	Issueid int64
	Content string
	Time    int64
}
type NoteByTimeDesc []Note

func (this NoteByTimeDesc) Len() int {
	return len(this)
}
func (this NoteByTimeDesc) Less(i, j int) bool {
	return this[i].Time > this[j].Time
}
func (this NoteByTimeDesc) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type Response struct {
	Success bool
	Message string
	Id      int64
}

var projects = []Project{}
var issues = []Issue{}
var notes = []Note{}

func loadProjects() error {
	filename := dataPath + "project.json"
	b, _ := ioutil.ReadFile(filename)

	err := json.Unmarshal(b, &projects)

	return err
}

func loadIssues() error {
	filename := dataPath + "issue.json"
	b, _ := ioutil.ReadFile(filename)

	err := json.Unmarshal(b, &issues)

	return err
}

func loadNotes() error {
	filename := dataPath + "note.json"
	b, _ := ioutil.ReadFile(filename)

	err := json.Unmarshal(b, &notes)

	return err
}

// exists returns whether the given file or directory exists or not
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func loadData() {
	err := loadProjects()
	err2 := loadIssues()
	err3 := loadNotes()

	if err != nil {
		log.Fatal("loadProjects: Unable to load. Was the data file removed? ")
	}

	if err2 != nil {
		log.Fatal("loadIssues: Unable to load. Was the data file removed? ")
	}

	if err3 != nil {
		log.Fatal("loadNotes: Unable to load. Was the data file removed? ")
	}

}

func writeProjects() error {
	filename := dataPath + "project.json"
	t, _ := json.Marshal(projects)

	err := ioutil.WriteFile(filename, t, 0600)

	return err
}

func writeIssues() error {
	filename := dataPath + "issue.json"
	t, _ := json.Marshal(issues)

	err := ioutil.WriteFile(filename, t, 0600)

	return err
}

func writeNotes() error {
	filename := dataPath + "note.json"
	t, _ := json.Marshal(notes)

	err := ioutil.WriteFile(filename, t, 0600)

	return err
}

func writeData() {
	err := writeProjects()
	err2 := writeIssues()
	err3 := writeNotes()

	if err != nil {
		log.Fatal("writeProjects: Unable to write. ", err)
	}

	if err2 != nil {
		log.Fatal("writeIssues: Unable to write. ", err2)
	}

	if err3 != nil {
		log.Fatal("writeNotes: Unable to write. ", err3)
	}
}

func setup() {
	// If the data directory does not exist lets create it
	// and then populate our data
	_, err := os.Stat(dataPath)
	if err != nil {
		mkerr := os.Mkdir(dataPath, 0700)

		if mkerr != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Was not able to create directory 'data' in current path.\n")
			os.Exit(2)
		}
		writeData()
	}

	loadData()
}

var templates = template.Must(template.ParseFiles(tmplPath+"base.html", tmplPath+"styles.css"))

// View handlers below
func baseHandler(w http.ResponseWriter, r *http.Request) {

	var loadtmpl = "base.html"
	templates.ExecuteTemplate(w, loadtmpl, "")

	/*if *useTemplates == true {
		body, err := ioutil.ReadFile(tmplPath + "base.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Was not able to read template base.html\n")
			os.Exit(2)
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, string(body))

	} else {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, basetemplate)
	}*/
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	var loadtmpl = "styles.css"
	w.Header().Set("Content-Type", "text/css")
	templates.ExecuteTemplate(w, loadtmpl, "")
}

// Returns all projects
func allprojectsHandler(w http.ResponseWriter, r *http.Request) {

	includeclosed, err := strconv.ParseBool(r.FormValue("c"))
	
	if err != nil {
		includeclosed = false
	}

	var proj = []Project{}
	for i := range projects {
	
		if includeclosed == true || projects[i].Open == true {
			proj = append(proj, projects[i])
		}
	}

	sort.Sort(ProjectByTimeDesc(proj))

	var projLimit = []Project{}
	if len(proj) < returnLimit {
		projLimit = proj
	} else {
		projLimit = proj[:returnLimit]
	}
	
	t, _ := json.Marshal(projLimit)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Returns all issues
func allissuesHandler(w http.ResponseWriter, r *http.Request) {

	includeclosed, err := strconv.ParseBool(r.FormValue("c"))
	
	if err != nil {
		includeclosed = false
	}

	var iss = []Issue{}
	for i := range issues {
		if includeclosed == true || issues[i].Open == true {
			iss = append(iss, issues[i])
		}
	}

	sort.Sort(IssueByTimeDesc(iss))

	var issLimit = []Issue{}
	if len(iss) < returnLimit {
		issLimit = iss
	} else {
		issLimit = iss[:returnLimit]
	}
	
	t, _ := json.Marshal(issLimit)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Returns all issues for supplied projectid
func issuesbyprojectHandler(w http.ResponseWriter, r *http.Request) {
	var iss = []Issue{}
	
	includeclosed, err := strconv.ParseBool(r.FormValue("c"))
	
	if err != nil {
		includeclosed = false
	}
	
	id, err := strconv.ParseInt(r.FormValue("q"), 10, 64)

	if err == nil {
		for i := range issues {
			if issues[i].Projectid == id && (includeclosed == true || issues[i].Open == true) {
				iss = append(iss, issues[i])
			}
		}
		sort.Sort(IssueByTimeDesc(iss))
	}

	t, _ := json.Marshal(iss)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Returns all notes for the supplied issueid
func notesbyissueHandler(w http.ResponseWriter, r *http.Request) {
	var iss = []Note{}
	id, err := strconv.ParseInt(r.FormValue("q"), 10, 64)

	if err == nil {
		for i := range notes {
			if notes[i].Issueid == id {
				iss = append(iss, notes[i])
			}
		}
		sort.Sort(NoteByTimeDesc(iss))
	}

	t, _ := json.Marshal(iss)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// This also needs to search notes inside the issues to be really
// effective
// by default this is an AND search
func issuessearchHandler(w http.ResponseWriter, r *http.Request) {
	var iss = []Issue{}
	term := r.FormValue("q")
	
	includeclosed, err := strconv.ParseBool(r.FormValue("c"))
	
	if err != nil {
		includeclosed = false
	}

	var words = strings.Split(strings.ToLower(term), " ")

	for i := range issues {
		match := true

		noteText := ""
		// Get all of the notes for this issue
		for j := range notes {
			if notes[j].Issueid == issues[i].Id {
				noteText = noteText + notes[j].Content
			}
		}

		// Get the project name for this issue
		projectName := ""
		for k := range projects {
			if projects[k].Id == issues[i].Projectid && (includeclosed == true || projects[k].Open == true) {
				projectName = projects[k].Name
			}
		}

		for _, word := range words {
			check := strings.ToLower(fmt.Sprintf("%s %s %s #%d", projectName, noteText, issues[i].Name, issues[i].Id))

			if strings.Contains(check, word) == false {
				match = false
			}
		}

		if match == true && (includeclosed == true || issues[i].Open == true) {
			iss = append(iss, issues[i])
		}
	}

	sort.Sort(IssueByTimeDesc(iss))

	var issLimit = []Issue{}
	if len(iss) < returnLimit {
		issLimit = iss
	} else {
		issLimit = iss[:returnLimit]
	}
	
	t, _ := json.Marshal(issLimit)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Search over the projects OR style
// The reason being we want to display any name thats even
// close to being what we want
func projectssearchHandler(w http.ResponseWriter, r *http.Request) {
	var proj = []Project{}
	term := r.FormValue("q")
	
	includeclosed, err := strconv.ParseBool(r.FormValue("c"))
	
	if err != nil {
		includeclosed = false
	}

	var words = strings.Split(strings.ToLower(term), " ")

	// for each project
	// get all of its issues
	// and all of that issues notes
	for i := range projects {
		if includeclosed == true || projects[i].Open == true {
			projectid := projects[i].Id
			check := projects[i].Name
			match := false

			for j := range issues {
				if issues[j].Projectid == projectid && (includeclosed == true || issues[j].Open == true) {
					issueid := issues[j].Id
					check = check + fmt.Sprintf("%s #%d", issues[j].Name, issues[j].Id)

					for k := range notes {

						if notes[k].Issueid == issueid && issues[j].Open == true {
							check = check + notes[k].Content
						}
					}

				}
			}
			check = strings.ToLower(check)

			for _, word := range words {
				if strings.Contains(check, word) == true {
					match = true
				}
			}

			if match == true {
				proj = append(proj, projects[i])
			}
		}
	}

	sort.Sort(ProjectByTimeDesc(proj))

	
	var projLimit = []Project{}
	if len(proj) < returnLimit {
		projLimit = proj
	} else {
		projLimit = proj[:returnLimit]
	}
	
	t, _ := json.Marshal(projLimit)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Generic method for saving a note
// it generates the new ID itself 
func savenote(n Note) {
	var topid int64
	topid = -1
	for i := range notes {
		if notes[i].Id > topid {
			topid = notes[i].Id
		}
	}
	topid = topid + 1

	n.Id = topid

	notes = append(notes, n)
	go writeNotes()
}

// Saves a note when supplied valid details
// also updates its issue so it pops to the top of searches
func savenoteHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{true, "", 0}
	issueid, err := strconv.ParseInt(r.FormValue("issueid"), 10, 64)
	body := html.EscapeString(r.FormValue("body"))

	// Validations
	if strings.Trim(body, " ") == "" {
		response.Success = false
		response.Message = "Content must not be empty."
	}

	if err != nil {
		response.Success = false
		response.Message = "IssueID must be integer."
	}

	if response.Success == true {

		// Get the related issue and update its time
		// also set it to be open
		for i := range issues {
			if issues[i].Id == issueid {
				issues[i].Time = int64(time.Now().Unix())
				issues[i].Open = true
			}
		}

		savenote(Note{0, issueid, body, int64(time.Now().Unix())})
		go writeIssues()
	}

	t, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Saves an issue when supplied valid details
// gets new ID based on the current highest ID
func saveissueHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{true, "", 0}
	projectid, err := strconv.ParseInt(r.FormValue("projectid"), 10, 64)
	issuename := html.EscapeString(r.FormValue("issuename"))
	issuecontent := html.EscapeString(r.FormValue("issuecontent"))

	// Validations
	if strings.Trim(issuename, " ") == "" {
		response.Success = false
		response.Message = "Issue name must not be empty."
	}

	if strings.Trim(issuecontent, " ") == "" {
		response.Success = false
		response.Message = "Issue content must not be empty."
	}

	if err != nil {
		response.Success = false
		response.Message = "ProjectID must be integer."
	}

	// If we are the demo version check we dont have 5 issues already
	if demoVersion == true && len(issues) >= 5 {
		response.Success = false
		response.Message = "Demo version is limited to 5 issues."
	}

	if response.Success == true {
		var topid int64
		topid = -1
		for i := range issues {
			if issues[i].Id > topid {
				topid = issues[i].Id
			}
		}
		topid = topid + 1

		issues = append(issues, Issue{topid, projectid, issuename, true, int64(time.Now().Unix())})

		response.Id = topid

		// Save the note as well using the ID we have here
		savenote(Note{0, topid, issuecontent, int64(time.Now().Unix())})

		go writeIssues()
	}

	t, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Saves a project when supplied a name
// gets new ID based on the current highest ID
func saveprojectHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{true, "", 0}
	projectname := html.EscapeString(r.FormValue("projectname"))

	// Validations
	if strings.Trim(projectname, " ") == "" {
		response.Success = false
		response.Message = "Project Name must not be empty."
	}

	for i := range projects {
		if projects[i].Name == projectname && projects[i].Open == true {
			response.Success = false
			response.Message = "Project Name must not already exist."
		}
	}

	// If we are the demo version check we dont have 1 project already
	if demoVersion == true && len(projects) >= 1 {
		response.Success = false
		response.Message = "Demo version is limited to 1 project."
	}

	// If we have passed validation actually save it
	if response.Success == true {
		var topid int64
		topid = -1
		for i := range projects {
			if projects[i].Id > topid {
				topid = projects[i].Id
			}
		}
		topid = topid + 1

		projects = append(projects, Project{topid, projectname, true, int64(time.Now().Unix())})
		go writeProjects()
	}

	t, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Closes down issues by setting the appropiate id to closed
func closeissueHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{true, "", 0}
	issueid, err := strconv.ParseInt(r.FormValue("issueid"), 10, 64)

	// Validations
	if err != nil {
		response.Success = false
		response.Message = "IssueID must be integer."
	}

	if response.Success == true {

		for i := range issues {
			if issues[i].Id == issueid {
				issues[i].Open = false
			}
		}

		go writeIssues()
	}

	t, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// Closes down projects by setting it to closed and
// all of its childen issues
func closeprojectHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{true, "", 0}
	projectid, err := strconv.ParseInt(r.FormValue("projectid"), 10, 64)

	// Validations
	if err != nil {
		response.Success = false
		response.Message = "ProjectID must be integer."
	}

	if response.Success == true {

		for i := range projects {
			if projects[i].Id == projectid {
				projects[i].Open = false
			}
		}
		
		for i := range issues {
			if issues[i].Projectid == projectid {
				issues[i].Open = false
			}
		}

		go writeProjects()
		go writeIssues()
	}

	t, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(t))
}

// TODO investigate if this is broken
// Should display what ips we are listening on but appears to have
// issues as it always displays things like 127.0.1.1 which seems
// odd.... NB the extra 1
func displayListening() {
	name, err := os.Hostname()
	if err == nil {

		addrs, err := net.LookupHost(name)

		if err == nil {
			for _, a := range addrs {
				log.Println("Listening on " + a)
			}
		}
	}
}

func main() {
	flag.Parse()

	setup()
	//displayListening()

	http.HandleFunc("/allprojects/", allprojectsHandler)
	http.HandleFunc("/allissues/", allissuesHandler)
	http.HandleFunc("/issuesbyproject/", issuesbyprojectHandler)
	http.HandleFunc("/notesbyissue/", notesbyissueHandler)
	http.HandleFunc("/issuessearch/", issuessearchHandler)
	http.HandleFunc("/projectssearch/", projectssearchHandler)
	http.HandleFunc("/savenote/", savenoteHandler)
	http.HandleFunc("/saveissue/", saveissueHandler)
	http.HandleFunc("/saveproject/", saveprojectHandler)
	http.HandleFunc("/closeissue/", closeissueHandler)
	http.HandleFunc("/closeproject/", closeprojectHandler)
	http.HandleFunc("/css/", cssHandler)
	http.HandleFunc("/", baseHandler)

	log.Println(fmt.Sprintf("Ready to serve requests on port %d", *portNumber))

	http.ListenAndServe(fmt.Sprintf(":%d", *portNumber), nil)
}
