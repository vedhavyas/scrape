package scrape

import "os"

// generateSiteMap will write the crawled url to given file
func generateSiteMap(fileName string, urls map[string]int) error {
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
	for loc := range urls {
		fh.WriteString("    " + "<url>\n")
		fh.WriteString("      " + "<loc>" + loc + "</loc>\n")
		fh.WriteString("      " + "<changefreq>weekly</changefreq>\n")
		fh.WriteString("      " + "<priority>0.5</priority>\n")
		fh.WriteString("    " + "</url>\n")
	}
	fh.WriteString("</urlset> ")

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
