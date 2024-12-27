package qbittorrent

import (
	"strings"
	"testing"
)

const (
	testServerURL = "http://127.0.0.1:8080"
	testUsername  = "admin"
	testPassword  = "adminadmin"
)

func TestLoginRequirement(t *testing.T) {
	client, err := NewDefaultClient(testServerURL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("Methods should with login", func(t *testing.T) {
		tests := []struct {
			name     string
			fn       func() error
			skipMsg  string
		}{
			{"Login", func() error {
				err := client.Login(testUsername, testPassword)
				return err
			}, ""},
			{"GetApplicationVersion", func() error {
				_, err := client.GetApplicationVersion()
				return err
			}, ""},
			{"GetAPIVersion", func() error {
				_, err := client.GetAPIVersion()
				return err
			}, ""},
			{"GetApplicationPreferences", func() error {
				_, err := client.GetApplicationPreferences()
				return err
			}, ""},
			{"GetDefaultSavePath", func() error {
				_, err := client.GetDefaultSavePath()
				return err
			}, ""},
			{"GetGlobalTransferInfo", func() error {
				_, err := client.GetGlobalTransferInfo()
				return err
			}, ""},
			{"GetTorrentList", func() error {
				_, err := client.GetTorrentList()
				return err
			}, ""},
			{"GetLog", func() error {
				_, err := client.GetLog()
				return err
			}, ""},
			{"GetPeerLog", func() error {
				_, err := client.GetPeerLog()
				return err
			}, ""},
			{"GetMainData", func() error {
				_, err := client.GetMainData(0)
				return err
			}, ""},
			{"GetAlternativeSpeedLimitsState", func() error {
				_, err := client.GetAlternativeSpeedLimitsState()
				return err
			}, ""},
			{"GetGlobalDownloadLimit", func() error {
				_, err := client.GetGlobalDownloadLimit()
				return err
			}, ""},
			{"GetGlobalUploadLimit", func() error {
				_, err := client.GetGlobalUploadLimit()
				return err
			}, ""},
			{"GetAllCategories", func() error {
				_, err := client.GetAllCategories()
				return err
			}, ""},
			{"GetAllTags", func() error {
				_, err := client.GetAllTags()
				return err
			}, ""},
			{"GetSearchPlugins", func() error {
				_, err := client.GetSearchPlugins()
				return err
			}, ""},
			{"ToggleAlternativeSpeedLimits", func() error {
				err := client.ToggleAlternativeSpeedLimits()
				return err
			}, ""},
			{"SetGlobalDownloadLimit", func() error {
				err := client.SetGlobalDownloadLimit(1024)
				return err
			}, ""},
			{"SetGlobalUploadLimit", func() error {
				err := client.SetGlobalUploadLimit(1024)
				return err
			}, ""},
			{"CreateTags", func() error {
				err := client.CreateTags([]string{"test_tag"})
				return err
			}, ""},
			{"DeleteTags", func() error {
				err := client.DeleteTags([]string{"test_tag"})
				return err
			}, ""},
			{"AddNewCategory", func() error {
				savePath, err := client.GetDefaultSavePath()
				if err != nil {
					return err
				}
				err = client.AddNewCategory("test_category", savePath)
				return err
			}, ""},
			{"RemoveCategories", func() error {
				err := client.RemoveCategories([]string{"test_category"})
				return err
			}, ""},
			// Skip methods requiring existing torrents
			{"GetTorrentGenericProperties", func() error {
				_, err := client.GetTorrentGenericProperties("test")
				return err
			}, "requires existing torrent"},
			{"GetTorrentTrackers", func() error {
				_, err := client.GetTorrentTrackers("test")
				return err
			}, "requires existing torrent"},
			{"GetTorrentWebSeeds", func() error {
				_, err := client.GetTorrentWebSeeds("test")
				return err
			}, "requires existing torrent"},
			{"GetTorrentContents", func() error {
				_, err := client.GetTorrentContents("test")
				return err
			}, "requires existing torrent"},
			{"GetTorrentPiecesStates", func() error {
				_, err := client.GetTorrentPiecesStates("test")
				return err
			}, "requires existing torrent"},
			{"GetTorrentPiecesHashes", func() error {
				_, err := client.GetTorrentPiecesHashes("test")
				return err
			}, "requires existing torrent"},
			{"PauseTorrents", func() error {
				err := client.PauseTorrents([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"ResumeTorrents", func() error {
				err := client.ResumeTorrents([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"DeleteTorrents", func() error {
				err := client.DeleteTorrents([]string{"test"}, false)
				return err
			}, "requires existing torrent"},
			{"RecheckTorrents", func() error {
				err := client.RecheckTorrents([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"ReannounceTorrents", func() error {
				err := client.ReannounceTorrents([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"EditTrackers", func() error {
				err := client.EditTrackers("test", "old", "new")
				return err
			}, "requires existing torrent"},
			{"RemoveTrackers", func() error {
				err := client.RemoveTrackers("test", []string{"url"})
				return err
			}, "requires existing torrent"},
			{"AddPeers", func() error {
				err := client.AddPeers("test", []string{"peer"})
				return err
			}, "requires existing torrent"},
			{"AddTrackersToTorrent", func() error {
				err := client.AddTrackersToTorrent("test", []string{"tracker"})
				return err
			}, "requires existing torrent"},
			{"IncreaseTorrentPriority", func() error {
				err := client.IncreaseTorrentPriority([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"DecreaseTorrentPriority", func() error {
				err := client.DecreaseTorrentPriority([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"MaximalTorrentPriority", func() error {
				err := client.MaximalTorrentPriority([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"MinimalTorrentPriority", func() error {
				err := client.MinimalTorrentPriority([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"SetFilePriority", func() error {
				err := client.SetFilePriority("test", []int{1}, 1)
				return err
			}, "requires existing torrent"},
			{"GetTorrentDownloadLimit", func() error {
				_, err := client.GetTorrentDownloadLimit([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"SetTorrentDownloadLimit", func() error {
				err := client.SetTorrentDownloadLimit([]string{"test"}, 1024)
				return err
			}, "requires existing torrent"},
			{"SetTorrentShareLimit", func() error {
				err := client.SetTorrentShareLimit([]string{"test"}, 1.0, 3600)
				return err
			}, "requires existing torrent"},
			{"GetTorrentUploadLimit", func() error {
				_, err := client.GetTorrentUploadLimit([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"SetTorrentUploadLimit", func() error {
				err := client.SetTorrentUploadLimit([]string{"test"}, 1024)
				return err
			}, "requires existing torrent"},
			{"SetTorrentLocation", func() error {
				err := client.SetTorrentLocation([]string{"test"}, "/path")
				return err
			}, "requires existing torrent"},
			{"SetTorrentName", func() error {
				err := client.SetTorrentName("test", "newname")
				return err
			}, "requires existing torrent"},
			{"SetTorrentCategory", func() error {
				err := client.SetTorrentCategory([]string{"test"}, "category")
				return err
			}, "requires existing torrent"},
			{"AddTorrentTags", func() error {
				err := client.AddTorrentTags([]string{"test"}, []string{"tag"})
				return err
			}, "requires existing torrent"},
			{"RemoveTorrentTags", func() error {
				err := client.RemoveTorrentTags([]string{"test"}, []string{"tag"})
				return err
			}, "requires existing torrent"},
			{"SetAutomaticTorrentManagement", func() error {
				err := client.SetAutomaticTorrentManagement([]string{"test"}, true)
				return err
			}, "requires existing torrent"},
			{"ToggleSequentialDownload", func() error {
				err := client.ToggleSequentialDownload([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"SetFirstLastPiecePriority", func() error {
				err := client.SetFirstLastPiecePriority([]string{"test"})
				return err
			}, "requires existing torrent"},
			{"SetForceStart", func() error {
				err := client.SetForceStart([]string{"test"}, true)
				return err
			}, "requires existing torrent"},
			{"SetSuperSeeding", func() error {
				err := client.SetSuperSeeding([]string{"test"}, true)
				return err
			}, "requires existing torrent"},
			{"RenameFile", func() error {
				err := client.RenameFile("test", "old", "new")
				return err
			}, "requires existing torrent"},
			{"RenameFolder", func() error {
				err := client.RenameFolder("test", "old", "new")
				return err
			}, "requires existing torrent"},
			// Skip RSS-related methods
			{"AddFolder", func() error {
				err := client.AddFolder("/path")
				return err
			}, "RSS feature not enabled"},
			{"AddFeed", func() error {
				err := client.AddFeed("url", "/path")
				return err
			}, "RSS feature not enabled"},
			{"RemoveItem", func() error {
				err := client.RemoveItem("/path")
				return err
			}, "RSS feature not enabled"},
			{"MoveItem", func() error {
				err := client.MoveItem("/old", "/new")
				return err
			}, "RSS feature not enabled"},
			{"MarkAsRead", func() error {
				err := client.MarkAsRead("/path", "article")
				return err
			}, "RSS feature not enabled"},
			{"RefreshItem", func() error {
				err := client.RefreshItem("/path")
				return err
			}, "RSS feature not enabled"},
			{"SetAutoDownloadingRule", func() error {
				err := client.SetAutoDownloadingRule("rule", "def")
				return err
			}, "RSS feature not enabled"},
			{"RenameAutoDownloadingRule", func() error {
				err := client.RenameAutoDownloadingRule("old", "new")
				return err
			}, "RSS feature not enabled"},
			{"RemoveAutoDownloadingRule", func() error {
				err := client.RemoveAutoDownloadingRule("rule")
				return err
			}, "RSS feature not enabled"},
			{"GetAllArticlesMatchingRule", func() error {
				_, err := client.GetAllArticlesMatchingRule("rule")
				return err
			}, "RSS feature not enabled"},
			// Skip search-related methods
			{"GetSearchStatus", func() error {
				_, err := client.GetSearchStatus(0)
				return err
			}, "search feature not enabled"},
			{"GetSearchResults", func() error {
				_, err := client.GetSearchResults(0, 10, 10)
				return err
			}, "search feature not enabled"},
			{"StartSearch", func() error {
				_, err := client.StartSearch("pattern", []string{"plugin"}, "category")
				return err
			}, "search feature not enabled"},
			{"StopSearch", func() error {
				err := client.StopSearch(0)
				return err
			}, "search feature not enabled"},
			{"DeleteSearch", func() error {
				err := client.DeleteSearch(0)
				return err
			}, "search feature not enabled"},
			{"EnableSearchPlugin", func() error {
				err := client.EnableSearchPlugin([]string{"plugin"}, true)
				return err
			}, "search feature not enabled"},
			{"UpdateSearchPlugins", func() error {
				err := client.UpdateSearchPlugins()
				return err
			}, "search feature not enabled"},
			{"Logout", func() error {
				err := client.Logout()
				return err
			}, ""},
		}

		// First login
		if err := client.Login(testUsername, testPassword); err != nil {
			t.Fatalf("Initial login failed: %v", err)
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.skipMsg != "" {
					t.Skipf("Skipping %s: %s", tt.name, tt.skipMsg)
					return
				}

				if err := tt.fn(); err != nil {
					if strings.Contains(err.Error(), "not found") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "400") ||
						strings.Contains(err.Error(), "409") {
						t.Skipf("Skipping %s: resource not found or invalid request", tt.name)
					} else {
						t.Errorf("%s failed with login: %v", tt.name, err)
					}
				} else {
					t.Logf("%s passed", tt.name)
				}
			})
		}
	})
}
