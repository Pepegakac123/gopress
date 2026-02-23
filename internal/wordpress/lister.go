package wordpress

import (
	"fmt"
	"strings"
)

// ListImagesFromFolder szuka folderu o podanej nazwie i wypisuje linki do zdjęć w podfolderach.
// Zwraca mapę: nazwa podfolderu -> lista linków.
func (c *Client) ListImagesFromFolder(rootName, rootNameShort string) (map[string][]string, error) {
	folders, err := c.GetFileBirdFolders()
	if err != nil {
		return nil, err
	}

	root := findFolderRecursive(folders, rootName, rootNameShort)
	if root == nil {
		return nil, fmt.Errorf("nie znaleziono folderu o nazwie zawierającej \"%s\" lub \"%s\"", rootName, rootNameShort)
	}

	results := make(map[string][]string)

	for _, sub := range root.Children {
		attachmentIDs, err := c.GetFileBirdAttachmentIDs(sub.ID)
		if err != nil {
			// Możemy tu zalogować błąd lub go zignorować dla danego podfolderu
			continue
		}

		links := make([]string, 0, len(attachmentIDs))
		for _, id := range attachmentIDs {
			url, err := c.GetMediaSourceURL(id)
			if err != nil {
				continue
			}
			links = append(links, url)
		}
		results[sub.Text] = links
	}

	return results, nil
}

func findFolderRecursive(folders []FbFolder, name, shortName string) *FbFolder {
	n := strings.ToLower(name)
	ns := strings.ToLower(shortName)

	for _, f := range folders {
		lowerText := strings.ToLower(f.Text)
		if strings.Contains(lowerText, n) || (ns != "" && strings.Contains(lowerText, ns)) {
			return &f
		}
		if len(f.Children) > 0 {
			res := findFolderRecursive(f.Children, name, shortName)
			if res != nil {
				return res
			}
		}
	}
	return nil
}
