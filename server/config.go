package coditra

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fauu/coditra/lookup"

	ntgo "github.com/dolow/nt-go"
)

const (
	configFileName   = "config.nt"
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
)

//go:embed config.sample.nt
var sampleConfig string

// Config is the representations of the program's configuration in its memory
type Config struct {
	DocumentsDir string
	UserAgent    string
	Setups       []LookupSetup
}

type LookupSetup struct {
	Name          string        `json:"name"`
	LookupEntries []LookupEntry `json:"lookupEntries"`
}

type LookupEntry struct {
	ID             string               `nt:"id" json:"id"`
	Name           *string              `nt:"name" json:"name"`
	LookupSourceID *string              `json:"source"`
	URL            *string              `nt:"url" json:"url"`
	Params         *lookup.LookupParams `nt:"params" json:"params"`
}

// FSConfig is the representation the program's configuration in the user's filesystem
type FSConfig struct {
	DocumentsDir  string          `nt:"documentsDir"`
	UserAgent     *string         `nt:"userAgent"`
	LookupEntries []LookupEntry   `nt:"lookups"`
	LookupSetups  []FSLookupSetup `nt:"setups"`
}

type FSLookupSetup struct {
	Name            string   `nt:"name"`
	LookupEntryRefs []string `nt:"lookups"` // Each parses into FSLookupEntryRef
}

type FSLookupEntryRef struct {
	LookupSourceID string
	Params         *lookup.LookupParams
}

func LoadConfig() (*Config, error) {
	configDirPath, err := GetUserConfigLocation(AppName)
	if err != nil {
		return nil, fmt.Errorf("getting user config location: %v", err)
	}

	filePath := path.Join(configDirPath, configFileName)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Cannot read config file (%v). Trying to create sample config…", err)

		err := ensureConfigDir(configDirPath)
		if err != nil {
			return nil, fmt.Errorf("creating config directory: %v", err)
		}

		err = createSampleConfig(filePath)
		if err != nil {
			return nil, fmt.Errorf("creating sample config file: %v", err)
		} else {
			fmt.Printf(
				"Created sample config file at '%s'\n"+
					"Please specify the translation documents directory at"+
					" the top of that file and relaunch coditra\n",
				filePath,
			)
			return nil, errors.New("config file was not ready")
		}
	}

	fsConfig := FSConfig{}
	err = ntgo.Marshal(string(fileContent), &fsConfig)
	if err != nil {
		return nil, fmt.Errorf("parsing config file: %v", err)
	}

	if fsConfig.DocumentsDir == "" {
		return nil, errors.New("'documentsDir' not specified in config")
	}

	// nt-go can silently return empty slice if it encounters wrong syntax (e.g., missing colon after key name)
	if len(fsConfig.LookupEntries) == 0 {
		return nil, errors.New("no Lookups loaded from config file. If you have specified some, please check for errors in syntax")
	}

	if len(fsConfig.LookupSetups) == 0 {
		return nil, errors.New("no Setups loaded from the config file. If you have specified some, please check for errors in syntax")
	}

	for i, entry := range fsConfig.LookupEntries {
		if entry.ID == "" {
			return nil, errors.New("a Lookup needs to have an 'id'")
		}
		if (entry.LookupSourceID == nil) == (entry.URL == nil) {
			return nil, errors.New("a Lookup needs to have exactly one of 'source' and 'url' specified")
		}
		if entry.LookupSourceID != nil {
			for sourceID, source := range LookupSources {
				if sourceID == *entry.LookupSourceID {
					fsConfig.LookupEntries[i].Name = &source.Name
					break
				}
			}
		}
	}

	setups := []LookupSetup{}
	for _, fsSetup := range fsConfig.LookupSetups {
		if fsSetup.Name == "" {
			return nil, errors.New("a Setup needs to have a 'name'")
		}
		if len(fsSetup.LookupEntryRefs) == 0 {
			return nil, errors.New("a Setup needs to have at least one Lookup")
		}

		entries := []LookupEntry{}
	entryRefs:
		for _, entryRef := range fsSetup.LookupEntryRefs {
			if entryRef == "" {
				return nil, errors.New("a Lookup in a Setup cannot be empty")
			}
			for _, entry := range fsConfig.LookupEntries {
				if entry.ID == entryRef {
					entries = append(entries, entry)
					continue entryRefs
				}
			}
			parsedEntryRef := parseFSLookupEntryRef(entryRef)
			for id, source := range LookupSources {
				if id == parsedEntryRef.LookupSourceID {
					params := parsedEntryRef.Params
					if params == nil {
						params = &source.ConstantParams
					}
					entry := LookupEntry{
						ID:             entryRef,
						LookupSourceID: &parsedEntryRef.LookupSourceID,
						URL:            nil,
						Name:           &source.Name,
						Params:         params,
					}
					entries = append(entries, entry)
					continue entryRefs
				}
			}
			return nil, fmt.Errorf("no corresponding Lookup found for ref '%s'", entryRef)
		}

		setups = append(setups, LookupSetup{
			Name:          fsSetup.Name,
			LookupEntries: entries,
		})
	}

	// IMPROVEMENT: Validate that sources have valid params right here

	userAgent := fsConfig.UserAgent
	if userAgent == nil || len(strings.TrimSpace(*userAgent)) == 0 {
		var defaultUserAgentCopy = defaultUserAgent
		userAgent = &defaultUserAgentCopy
	}

	return &Config{
		DocumentsDir: fsConfig.DocumentsDir,
		UserAgent:    *userAgent,
		Setups:       setups,
	}, nil
}

func parseFSLookupEntryRef(s string) FSLookupEntryRef {
	s = strings.Replace(s, "(", " ", 1)
	s = strings.Replace(s, ")", " ", 1)
	s = strings.Replace(s, ",", " ", 1)
	segs := strings.Fields(s)
	// ASSUMPTION: s is non-empty
	out := FSLookupEntryRef{LookupSourceID: segs[0]}
	l := len(segs)
	if l >= 2 {
		out.Params = &lookup.LookupParams{}
		out.Params.SourceLang = &segs[1]
	}
	if l >= 3 {
		out.Params.TargetLang = &segs[2]
	}
	return out
}

func ensureConfigDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("running mkdir: %v", err)
		}
	}
	return nil
}

func createSampleConfig(filePath string) error {
	err := os.WriteFile(filePath, []byte(sampleConfig), 0644)
	if err != nil {
		return fmt.Errorf("creating file: %v", err)
	}
	return nil
}
