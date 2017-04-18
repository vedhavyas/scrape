package utils

import "os"

//GenerateSiteMap will write the crawled url to given file
func GenerateSiteMap(fileName string, urls []string) error {

	err := deleteFileIfExists(fileName)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fh.Close()

	fh.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fh.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, loc := range urls {
		fh.WriteString("    " + "<url>\n")
		fh.WriteString("      " + "<loc>" + loc + "</loc>\n")
		fh.WriteString("      " + "<changefreq>weekly</changefreq>\n")
		fh.WriteString("      " + "<priority>0.5</priority>\n")
		fh.WriteString("    " + "</url>\n")
	}
	fh.WriteString("</urlset> ")

	return nil

}

//GenerateAssetFile writes assets fetched from each page to given file
func GenerateAssetFile(fileName string, assetInPages map[string][]string) error {
	err := deleteFileIfExists(fileName)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fh.Close()

	for pageURL, assets := range assetInPages {
		if len(assets) == 0 {
			continue
		}
		fh.WriteString(pageURL + "\n")
		for _, assetURL := range assets {
			fh.WriteString("    " + " - " + assetURL + "\n")
		}
	}

	return nil

}

//deleteFileIfExists deletes a file if exists
func deleteFileIfExists(fileName string) error {
	//delete old file first
	if _, err := os.Stat(fileName); err == nil {
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
