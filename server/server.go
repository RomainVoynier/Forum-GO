package server

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

type Topic struct {
    Title     string `json:"title"`
    Author    string `json:"author"`
    Content   string `json:"content"`
    CreatedAt string `json:"createdAt"`
}

var (
    topics   []Topic
    topicsMu sync.Mutex
)

const topicsFile = "data/topics.json"

func RegisterHandlers() {
    if err := loadTopics(); err != nil {
        log.Println("Warning loading topics:", err)
    }

    http.HandleFunc("/topics", topicsHandler)
    http.HandleFunc("/topic/add", addTopicHandler)
    http.HandleFunc("/topic/delete", deleteTopicHandler)
    http.HandleFunc("/create_post.html", createPostHandler)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    http.ServeFile(w, r, "html/create_post.html")
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    topicsMu.Lock()
    defer topicsMu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(topics)
}

func addTopicHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var topic Topic
    contentType := r.Header.Get("Content-Type")
    if strings.HasPrefix(contentType, "application/json") {
        if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
            http.Error(w, "Données JSON invalides", http.StatusBadRequest)
            return
        }
    } else {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "Impossible de lire le formulaire", http.StatusBadRequest)
            return
        }
        topic.Title = r.FormValue("title")
        topic.Author = r.FormValue("author")
        topic.Content = r.FormValue("content")
    }

    topic.Title = strings.TrimSpace(topic.Title)
    topic.Author = strings.TrimSpace(topic.Author)
    topic.Content = strings.TrimSpace(topic.Content)

    if topic.Title == "" {
        http.Error(w, "Le titre est requis", http.StatusBadRequest)
        return
    }

    if topic.Content == "" {
        http.Error(w, "Le contenu est requis", http.StatusBadRequest)
        return
    }

    if topic.Author == "" {
        topic.Author = "Anonyme"
    }

    topic.CreatedAt = time.Now().Format("02/01/2006 15:04")

    topicsMu.Lock()
    defer topicsMu.Unlock()

    topics = append([]Topic{topic}, topics...)
    if err := saveTopics(); err != nil {
        log.Println("Erreur en sauvegardant les sujets:", err)
        http.Error(w, "Impossible d'enregistrer le fil", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func loadTopics() error {
    dir := filepath.Dir(topicsFile)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    data, err := os.ReadFile(topicsFile)
    if err != nil {
        if os.IsNotExist(err) {
            topics = []Topic{}
            return saveTopics()
        }
        return err
    }

    if len(data) == 0 {
        topics = []Topic{}
        return nil
    }

    return json.Unmarshal(data, &topics)
}

func saveTopics() error {
    data, err := json.MarshalIndent(topics, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(topicsFile, data, 0644)
}

func deleteTopicHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    var payload struct {
        Title     string `json:"title"`
        CreatedAt string `json:"createdAt"`
    }

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Données JSON invalides", http.StatusBadRequest)
        return
    }

    payload.Title = strings.TrimSpace(payload.Title)
    payload.CreatedAt = strings.TrimSpace(payload.CreatedAt)
    if payload.Title == "" {
        http.Error(w, "Le titre est requis pour supprimer un fil", http.StatusBadRequest)
        return
    }

    topicsMu.Lock()
    defer topicsMu.Unlock()

    index := -1
    for i, topic := range topics {
        if topic.Title == payload.Title && (payload.CreatedAt == "" || topic.CreatedAt == payload.CreatedAt) {
            index = i
            break
        }
    }

    if index == -1 {
        http.Error(w, "Aucun fil trouvé avec ce titre", http.StatusNotFound)
        return
    }

    topics = append(topics[:index], topics[index+1:]...)
    if err := saveTopics(); err != nil {
        log.Println("Erreur en sauvegardant les sujets:", err)
        http.Error(w, "Impossible de supprimer le fil", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
