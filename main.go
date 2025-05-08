package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"schartz.com/sg/helpers"
	"schartz.com/sg/templates"
)


func md_to_html(md_file_path, html_file_path string) error {
	md_content, err := os.ReadFile(md_file_path)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	html_str := markdown.ToHTML(md_content, nil, nil)
	html_title:= helpers.MakeTitle(filepath.Base(md_file_path))
	html_str = fmt.Appendf(nil, templates.GetCommonHtmlTpl(), html_title, html_str)

	err = os.WriteFile(html_file_path, html_str, 0644)
	if err != nil {
		return fmt.Errorf("error trying to write filee to %s. Error is: %w", html_file_path, err)
	}
	return nil;

}

func convert_md_folder(in_folder_path, out_folder_path string) error {
	err := filepath.Walk(in_folder_path, func(in_md_file_path string, info os.FileInfo, err error) error {
		if err != nil  {
			return err
		}

		intermediate_out_path := filepath.Join(out_folder_path, strings.Trim(in_md_file_path, in_folder_path))

		if info.IsDir(){
			err := os.MkdirAll(intermediate_out_path, 0755)
			if err != nil {
				return fmt.Errorf("error creating dir %s. Error is %w", intermediate_out_path, err)
			}

			return nil
		}

		if strings.ToLower(filepath.Ext(in_md_file_path)) == ".md" || strings.ToLower(filepath.Ext(in_md_file_path))== "markdown" {
			out_html_path := intermediate_out_path + ".html"
			fmt.Printf("Converting %s to %s", in_md_file_path, out_html_path)
			err = md_to_html(in_md_file_path, out_html_path)
			if err != nil {
				return fmt.Errorf("error converting file: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	return nil
}


func generate_menu_html(out_dir string) (string, error) {
	var menu_html strings.Builder
	menu_html.WriteString(templates.GetMenuTpl())

	// recursive menu builder function
	var menuBuilderFunc func(dir string, level int) error
	menuBuilderFunc = func(dir string, level int) error {
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			name := file.Name()
			path := filepath.Join(dir, name)
			// Create a category/subcategory from dir
			if file.IsDir() {
				menu_html.WriteString(fmt.Sprintf(`<li class="menu-item has-children">\n  <a href="%s/index.html" class="menu-link">%s</a>\n    <ul class="menu-dropdown">\n`, name, name))
				if err := menuBuilderFunc(path, level+1); err != nil {
					return err
				}
				menu_html.WriteString("    </ul>\n  </li>\n")
			} else if strings.ToLower(filepath.Ext(name)) == ".html" && name != "index.html" {
				relPath, err := filepath.Rel(out_dir, path)
				if err != nil {
					return err
				}

				menu_html.WriteString(fmt.Sprintf(
					`<li class="menu-item"><a href="/%s" class="menu-link">%s</a></li>\n`, 
					relPath, 
					strings.TrimSuffix(name, ".html"),
				))
			}
		}

		return nil
	}

	// start building the menu from top level out dir
	if err := menuBuilderFunc(out_dir, 0); err != nil {
		return "", err
	}

	menu_html.WriteString(`</ul>\n</nav>\n</div>\n`)
	return menu_html.String(), nil
}


func generate_index_page(in_dir string, out_dir string) error {
	files, err := os.ReadDir(in_dir)
	if err != nil {
		return err
	}

	var html_content strings.Builder
	html_content.WriteString(fmt.Sprintf(templates.GetMenuIndexPageTpl(), filepath.Base(in_dir)))

	for _, file := range files {
		name := file.Name()
		path := filepath.Join(in_dir, name)

		if strings.ToLower(filepath.Ext(name)) == ".html" && name == "index.html" {
			relPath, err := filepath.Rel(out_dir, path)
			if err != nil {
				return nil
			}
			html_content.WriteString(fmt.Sprintf(`<li><a href="/%s">%s</a></li>\n`, relPath, strings.TrimSuffix(name, ".html")))
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input_dir> <output_dir>")
		return
	}

	in_dir := os.Args[1]
	out_dir := os.Args[2]

	dir_info, err := os.Stat(in_dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: Input directory '%s' does not exist.\n", in_dir)
		} else {
			fmt.Printf("Error: Could not stat input directory '%s': %v\n", in_dir, err)
		}
		return
	}
	if !dir_info.IsDir() {
		fmt.Printf("Error: Input path '%s' is not a directory.\n", in_dir)
		return
	}

	// Ensure that a fresh output dir is available
	err = os.Remove(out_dir)
	if err != nil {
		fmt.Println("Could not clean the output dir. Exiting.")
		return
	}

	outputInfo, err := os.Stat(out_dir)
	if err != nil {
		if os.IsNotExist(err) {
			// Attempt to create the output directory.
			err = os.MkdirAll(out_dir, 0755)
			if err != nil {
				fmt.Printf("Error: Could not create output directory '%s': %v\n", out_dir, err)
				return
			}
			fmt.Printf("Created output directory: %s\n", out_dir)
		} else {
			fmt.Printf("Error: Could not stat output directory '%s': %v\n", out_dir, err)
			return
		}
	} else if !outputInfo.IsDir() {
		fmt.Printf("Error: Output path '%s' is not a directory.\n", out_dir)
		return
	}

	// Convert the folder to HTML.
	err = convert_md_folder(in_dir, out_dir)
	if err != nil {
		// Use %+v for the stack trace
		fmt.Printf("Error: %v\n", err) 
		return
	}

	fmt.Println("Conversion complete.")

}
