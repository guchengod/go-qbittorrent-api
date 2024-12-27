package qbittorrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type QBittorrentClient struct {
	baseURL    string
    client     *http.Client
    cookie     *http.Cookie
}

func NewClient(baseURL string, httpClient *http.Client, cookie *http.Cookie) (*QBittorrentClient, error) {
    
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is empty")
	}

	if httpClient == nil {
        httpClient = &http.Client{}
    }
    if cookie == nil {
        cookie = &http.Cookie{}
    }
    
    return &QBittorrentClient{
        baseURL:    baseURL,
        client:     httpClient,
        cookie:     cookie,
    }, nil
}

func NewDefaultClient(baseURL string) (*QBittorrentClient, error) {
    return NewClient(baseURL, nil, nil)
}

func (q *QBittorrentClient) GetHttpClient() *http.Client {
	if q.client == nil {
		return &http.Client{}
	}
	return q.client
}

func (q *QBittorrentClient) GetCookie() *http.Cookie {
	if q.cookie == nil {
		return &http.Cookie{}
	}
	return q.cookie
}

func (q *QBittorrentClient) Login(username, password string) error {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/auth/login", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	q.cookie = resp.Cookies()[0]
	return nil
}

func (q *QBittorrentClient) Logout() error {
	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/auth/logout", nil)
	if err != nil {
		return err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status code: %d", resp.StatusCode)
	}

	q.cookie = nil
	return nil
}

func (q *QBittorrentClient) GetApplicationVersion() (string, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/app/version", nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (q *QBittorrentClient) GetAPIVersion() (string, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/app/webapiVersion", nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}


func (q *QBittorrentClient) GetApplicationPreferences() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/app/preferences", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var preferences map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&preferences)
	if err != nil {
		return nil, err
	}

	return preferences, nil
}

func (q *QBittorrentClient) SetApplicationPreferences(preferences map[string]interface{}) error {
	jsonData, err := json.Marshal(preferences)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/app/setPreferences", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set preferences failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetDefaultSavePath() (string, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/app/defaultSavePath", nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (q *QBittorrentClient) GetLog() ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/log/main", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var log []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&log)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (q *QBittorrentClient) GetPeerLog() ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/log/peers", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var peerLog []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&peerLog)
	if err != nil {
		return nil, err
	}

	return peerLog, nil
}

// Sync
func (q *QBittorrentClient) GetMainData(rid int) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/sync/maindata?rid=%d", q.baseURL, rid), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mainData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&mainData)
	if err != nil {
		return nil, err
	}

	return mainData, nil
}

func (q *QBittorrentClient) GetTorrentPeersData(hash string, rid int) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/sync/torrentPeers?hash=%s&rid=%d", q.baseURL, hash, rid), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var peersData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&peersData)
	if err != nil {
		return nil, err
	}

	return peersData, nil
}

// Transfer Info
func (q *QBittorrentClient) GetGlobalTransferInfo() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/transfer/info", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transferInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&transferInfo)
	if err != nil {
		return nil, err
	}

	return transferInfo, nil
}

func (q *QBittorrentClient) GetAlternativeSpeedLimitsState() (bool, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/transfer/speedLimitsMode", nil)
	if err != nil {
		return false, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var state int
	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		return false, err
	}

	return state == 1, nil
}

func (q *QBittorrentClient) ToggleAlternativeSpeedLimits() error {
	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/transfer/toggleSpeedLimitsMode", nil)
	if err != nil {
		return err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("toggle speed limits failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetGlobalDownloadLimit() (int, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/transfer/downloadLimit", nil)
	if err != nil {
		return 0, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var limit int
	err = json.NewDecoder(resp.Body).Decode(&limit)
	if err != nil {
		return 0, err
	}

	return limit, nil
}

func (q *QBittorrentClient) SetGlobalDownloadLimit(limit int) error {
	data := url.Values{}
	data.Set("limit", fmt.Sprintf("%d", limit))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/transfer/setDownloadLimit", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set download limit failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetGlobalUploadLimit() (int, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/transfer/uploadLimit", nil)
	if err != nil {
		return 0, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var limit int
	err = json.NewDecoder(resp.Body).Decode(&limit)
	if err != nil {
		return 0, err
	}

	return limit, nil
}

func (q *QBittorrentClient) SetGlobalUploadLimit(limit int) error {
	data := url.Values{}
	data.Set("limit", fmt.Sprintf("%d", limit))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/transfer/setUploadLimit", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set upload limit failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) BanPeers(peers string) error {
	data := url.Values{}
	data.Set("peers", peers)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/transfer/banPeers", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ban peers failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// Torrent Management
func (q *QBittorrentClient) GetTorrentList() ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/torrents/info", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var torrents []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&torrents)
	if err != nil {
		return nil, err
	}

	return torrents, nil
}

func (q *QBittorrentClient) GetTorrentGenericProperties(hash string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/torrents/properties?hash=%s", q.baseURL, hash), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var properties map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&properties)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

// Torrent Management (continued)
func (q *QBittorrentClient) GetTorrentTrackers(hash string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/torrents/trackers?hash=%s", q.baseURL, hash), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trackers []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&trackers)
	if err != nil {
		return nil, err
	}

	return trackers, nil
}

func (q *QBittorrentClient) GetTorrentWebSeeds(hash string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/torrents/webseeds?hash=%s", q.baseURL, hash), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var webSeeds []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&webSeeds)
	if err != nil {
		return nil, err
	}

	return webSeeds, nil
}

func (q *QBittorrentClient) GetTorrentContents(hash string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/torrents/files?hash=%s", q.baseURL, hash), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var contents []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&contents)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (q *QBittorrentClient) GetTorrentPiecesStates(hash string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/torrents/pieceStates?hash=%s", q.baseURL, hash), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var states []string
	err = json.NewDecoder(resp.Body).Decode(&states)
	if err != nil {
		return nil, err
	}

	return states, nil
}

func (q *QBittorrentClient) GetTorrentPiecesHashes(hash string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/torrents/pieceHashes?hash=%s", q.baseURL, hash), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var hashes []string
	err = json.NewDecoder(resp.Body).Decode(&hashes)
	if err != nil {
		return nil, err
	}

	return hashes, nil
}

func (q *QBittorrentClient) PauseTorrents(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/pause", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pause torrents failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) ResumeTorrents(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/resume", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resume torrents failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) DeleteTorrents(hashes []string, deleteFiles bool) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("deleteFiles", fmt.Sprintf("%t", deleteFiles))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/delete", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete torrents failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RecheckTorrents(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/recheck", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recheck torrents failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) ReannounceTorrents(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/reannounce", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reannounce torrents failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) EditTrackers(hash string, originalUrl string, newUrl string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("originalUrl", originalUrl)
	data.Set("newUrl", newUrl)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/editTracker", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("edit trackers failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RemoveTrackers(hash string, urls []string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("urls", strings.Join(urls, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/removeTrackers", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove trackers failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) AddPeers(hash string, peers []string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("peers", strings.Join(peers, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/addPeers", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add peers failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// Torrent Management (continued)
func (q *QBittorrentClient) AddNewTorrent(urls []string, options map[string]string) error {
	data := url.Values{}
	for key, value := range options {
		data.Set(key, value)
	}
	data.Set("urls", strings.Join(urls, "\n"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/add", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add new torrent failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) AddTrackersToTorrent(hash string, urls []string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("urls", strings.Join(urls, "\n"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/addTrackers", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add trackers to torrent failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) IncreaseTorrentPriority(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/increasePrio", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("increase torrent priority failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) DecreaseTorrentPriority(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/decreasePrio", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("decrease torrent priority failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) MaximalTorrentPriority(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/topPrio", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("maximal torrent priority failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) MinimalTorrentPriority(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/bottomPrio", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("minimal torrent priority failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetFilePriority(hash string, fileIds []int, priority int) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("ids", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(fileIds)), "|"), "[]"))
	data.Set("priority", fmt.Sprintf("%d", priority))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/filePrio", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set file priority failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetTorrentDownloadLimit(hashes []string) (map[string]int, error) {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/torrents/downloadLimit?"+data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var limits map[string]int
	err = json.NewDecoder(resp.Body).Decode(&limits)
	if err != nil {
		return nil, err
	}

	return limits, nil
}

func (q *QBittorrentClient) SetTorrentDownloadLimit(hashes []string, limit int) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("limit", fmt.Sprintf("%d", limit))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setDownloadLimit", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set torrent download limit failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetTorrentShareLimit(hashes []string, ratioLimit float64, seedingTimeLimit int) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("ratioLimit", fmt.Sprintf("%f", ratioLimit))
	data.Set("seedingTimeLimit", fmt.Sprintf("%d", seedingTimeLimit))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setShareLimits", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set torrent share limit failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetTorrentUploadLimit(hashes []string) (map[string]int, error) {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/torrents/uploadLimit?"+data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var limits map[string]int
	err = json.NewDecoder(resp.Body).Decode(&limits)
	if err != nil {
		return nil, err
	}

	return limits, nil
}

func (q *QBittorrentClient) SetTorrentUploadLimit(hashes []string, limit int) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("limit", fmt.Sprintf("%d", limit))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setUploadLimit", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set torrent upload limit failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetTorrentLocation(hashes []string, location string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("location", location)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setLocation", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set torrent location failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetTorrentName(hash string, name string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("name", name)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/rename", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set torrent name failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetTorrentCategory(hashes []string, category string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("category", category)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setCategory", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set torrent category failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetAllCategories() (map[string]map[string]int, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/torrents/categories", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var categories map[string]map[string]int
	err = json.NewDecoder(resp.Body).Decode(&categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (q *QBittorrentClient) AddNewCategory(category string, savePath string) error {
	data := url.Values{}
	data.Set("category", category)
	data.Set("savePath", savePath)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/createCategory", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add new category failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) EditCategory(category string, savePath string) error {
	data := url.Values{}
	data.Set("category", category)
	data.Set("savePath", savePath)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/editCategory", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("edit category failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RemoveCategories(categories []string) error {
	data := url.Values{}
	data.Set("categories", strings.Join(categories, "\n"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/removeCategories", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove categories failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) AddTorrentTags(hashes []string, tags []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("tags", strings.Join(tags, ","))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/addTags", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add torrent tags failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RemoveTorrentTags(hashes []string, tags []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("tags", strings.Join(tags, ","))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/removeTags", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove torrent tags failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetAllTags() ([]string, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/torrents/tags", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tags []string
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (q *QBittorrentClient) CreateTags(tags []string) error {
	data := url.Values{}
	data.Set("tags", strings.Join(tags, ","))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/createTags", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create tags failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) DeleteTags(tags []string) error {
	data := url.Values{}
	data.Set("tags", strings.Join(tags, ","))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/deleteTags", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete tags failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetAutomaticTorrentManagement(hashes []string, enable bool) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("enable", fmt.Sprintf("%t", enable))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setAutoManagement", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set automatic torrent management failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) ToggleSequentialDownload(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/toggleSequentialDownload", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("toggle sequential download failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetFirstLastPiecePriority(hashes []string) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setFirstLastPiecePrio", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set first/last piece priority failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetForceStart(hashes []string, enable bool) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("value", fmt.Sprintf("%t", enable))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setForceStart", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set force start failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetSuperSeeding(hashes []string, enable bool) error {
	data := url.Values{}
	data.Set("hashes", strings.Join(hashes, "|"))
	data.Set("value", fmt.Sprintf("%t", enable))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/setSuperSeeding", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set super seeding failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RenameFile(hash string, oldPath string, newPath string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("oldPath", oldPath)
	data.Set("newPath", newPath)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/renameFile", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rename file failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RenameFolder(hash string, oldPath string, newPath string) error {
	data := url.Values{}
	data.Set("hash", hash)
	data.Set("oldPath", oldPath)
	data.Set("newPath", newPath)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/torrents/renameFolder", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rename folder failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// RSS (Experimental)
func (q *QBittorrentClient) AddFolder(path string) error {
	data := url.Values{}
	data.Set("path", path)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/addFolder", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add folder failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) AddFeed(urlStr string, path string) error {
	data := url.Values{}
	data.Set("url", urlStr)
	data.Set("path", path)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/addFeed", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add feed failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RemoveItem(path string) error {
	data := url.Values{}
	data.Set("path", path)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/removeItem", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove item failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) MoveItem(itemPath string, destPath string) error {
	data := url.Values{}
	data.Set("itemPath", itemPath)
	data.Set("destPath", destPath)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/moveItem", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("move item failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetAllItems() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/rss/items", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (q *QBittorrentClient) MarkAsRead(itemPath string, articleId string) error {
	data := url.Values{}
	data.Set("itemPath", itemPath)
	data.Set("articleId", articleId)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/markAsRead", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mark as read failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RefreshItem(itemPath string) error {
	data := url.Values{}
	data.Set("itemPath", itemPath)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/refreshItem", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh item failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) SetAutoDownloadingRule(ruleName string, ruleDef string) error {
	data := url.Values{}
	data.Set("ruleName", ruleName)
	data.Set("ruleDef", ruleDef)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/setRule", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set auto-downloading rule failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RenameAutoDownloadingRule(ruleName string, newRuleName string) error {
	data := url.Values{}
	data.Set("ruleName", ruleName)
	data.Set("newRuleName", newRuleName)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/renameRule", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rename auto-downloading rule failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) RemoveAutoDownloadingRule(ruleName string) error {
	data := url.Values{}
	data.Set("ruleName", ruleName)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/rss/removeRule", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove auto-downloading rule failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetAllAutoDownloadingRules() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/rss/rules", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rules map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

func (q *QBittorrentClient) GetAllArticlesMatchingRule(ruleName string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/rss/matchingArticles?ruleName=%s", q.baseURL, ruleName), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var articles []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// Search
func (q *QBittorrentClient) StartSearch(pattern string, plugins []string, category string) (int, error) {
	data := url.Values{}
	data.Set("pattern", pattern)
	data.Set("plugins", strings.Join(plugins, "|"))
	data.Set("category", category)

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/search/start", strings.NewReader(data.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var id int
	err = json.NewDecoder(resp.Body).Decode(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (q *QBittorrentClient) StopSearch(id int) error {
	data := url.Values{}
	data.Set("id", fmt.Sprintf("%d", id))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/search/stop", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stop search failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetSearchStatus(id int) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/search/status?id=%d", q.baseURL, id), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (q *QBittorrentClient) GetSearchResults(id int, limit int, offset int) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/search/results?id=%d&limit=%d&offset=%d", q.baseURL, id, limit, offset), nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (q *QBittorrentClient) DeleteSearch(id int) error {
	data := url.Values{}
	data.Set("id", fmt.Sprintf("%d", id))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/search/delete", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete search failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) GetSearchPlugins() ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/search/plugins", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var plugins []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&plugins)
	if err != nil {
		return nil, err
	}

	return plugins, nil
}

func (q *QBittorrentClient) InstallSearchPlugin(sources []string) error {
	data := url.Values{}
	data.Set("sources", strings.Join(sources, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/search/installPlugin", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("install search plugin failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) UninstallSearchPlugin(names []string) error {
	data := url.Values{}
	data.Set("names", strings.Join(names, "|"))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/search/uninstallPlugin", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("uninstall search plugin failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) EnableSearchPlugin(names []string, enable bool) error {
	data := url.Values{}
	data.Set("names", strings.Join(names, "|"))
	data.Set("enable", fmt.Sprintf("%t", enable))

	req, err := http.NewRequest("POST", q.baseURL+"/api/v2/search/enablePlugin", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("enable search plugin failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (q *QBittorrentClient) UpdateSearchPlugins() error {
	req, err := http.NewRequest("GET", q.baseURL+"/api/v2/search/updatePlugins", nil)
	if err != nil {
		return err
	}
	req.AddCookie(q.cookie)

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update search plugins failed with status code: %d", resp.StatusCode)
	}

	return nil
}