package main

import (
	"path/filepath"
	"os"
	"os/exec"
	"fmt"
	"io"
    "time"
)

var SrcFolder string

func main() {
	if len(os.Args) > 1 {
		SrcFolder, _ = filepath.Abs(os.Args[1])
	} else {
		SrcFolder, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

    exPath, _ := filepath.Abs(os.Args[0])
    exName := filepath.Base(exPath)

    fmt.Println("> Scanning src folder...")
    srcFile, err := os.Open(SrcFolder)
    if err != nil {
        if srcFile != nil {
            srcFile.Close()
        }
        fmt.Println(">!> Warning! Unable to open the src folder for scanning! If you are trying to update the API, delete your src folder, re-create it, and start over.")
        fmt.Println(">!> If you are installing in a clean src folder, this should not affect your installation.")
        time.Sleep(2 * time.Second)
        return
    }
    names, err := srcFile.Readdirnames(-1)
    if err != nil && err != io.EOF {
        fmt.Println(">!> Warning! Unable to scan src folder! If you are trying to update the API, delete your src folder, re-create it, and start over.")
        fmt.Println(">!> If you are installing in a clean src folder, this should not affect your installation.")
        time.Sleep(2 * time.Second)
        return
    }
    nc := 0
    for _, name := range names {
        if name != exName {
            nc++
        }
    }
    srcFile.Close()

    if nc > 0 {
        fmt.Println("> Unclean folder detected! Attempting to update and not install...")
        Update()
        return
    }

	if err := SetupRepo(); err != nil {
		return
	}
	if err := StructureFolder(); err != nil {
		return
	}

	fmt.Println("> The CodeFrayAPI was installed successfully!")
}

/* The git commands:
	git init
	git remote add origin https://github.com/lg198/CodeFrayAPI.git
	git pull origin master
*/

func setupError(err error, command string) {
	fmt.Println(" >!> There was an error in installation: " + command + ". Please delete the src folder and recreate it.")
	fmt.Println("\t" + err.Error())
}

func SetupRepo() error {
	fmt.Println("> Setting up repository...")
	initCommand := exec.Command("git", "init")
	initCommand.Dir = SrcFolder
	if err := initCommand.Run(); err != nil {
		setupError(err, "git init")
		return err
	}
	remoteCommand := exec.Command("git", "remote", "add", "origin", "https://github.com/lg198/CodeFrayAPI.git")
	remoteCommand.Dir = SrcFolder
	if res, err := remoteCommand.CombinedOutput(); err != nil {
		fmt.Println(string(res))
		setupError(err, "git remote")
		return err
	}
	pullCommand := exec.Command("git", "pull", "origin", "master")
	pullCommand.Dir = SrcFolder
	if err := pullCommand.Run(); err != nil {
		setupError(err, "git pull")
		return err
	}
	return nil
}

func StructureFolder() error {
	fmt.Println("> Restructuring directory...")

	os.Remove(filepath.Join(SrcFolder, "LICENSE"))
	os.Remove(filepath.Join(SrcFolder, "README.md"))

	if err := CopyDirContents(filepath.Join(SrcFolder, "src"), SrcFolder); err != nil {
		setupError(err, "Copying directory contents")
		return err
	}

	if err := os.RemoveAll(filepath.Join(SrcFolder, "src")); err != nil {
		setupError(err, "Removing src folder")
		return err
	}
	return nil
}

/*
    Update command: git pull origin master
*/

func Update() {
    fmt.Println("> Pulling changes...")
    remoteCommand := exec.Command("git", "pull", "origin", "master")
    remoteCommand.Dir = SrcFolder
    if res, err := remoteCommand.CombinedOutput(); err != nil {
        fmt.Println(string(res))
        setupError(err, "git pull")
        return
    }

    fmt.Println("> The CodeFrayAPI was updated successfully!")
}

func CopyFile(source string, dest string) error {
     sourcefile, err := os.Open(source)
     if err != nil {
         return err
     }

     defer sourcefile.Close()

     destfile, err := os.Create(dest)
     if err != nil {
         return err
     }

     defer destfile.Close()

     _, err = io.Copy(destfile, sourcefile)
     if err == nil {
         sourceinfo, err := os.Stat(source)
         if err != nil {
             err = os.Chmod(dest, sourceinfo.Mode())
             if err != nil {
             	return err
             }
         }
     } else {
     	return err
     }

     return nil
 }

 func CopyDirContents(source string, dest string) error {
     directory, _ := os.Open(source)

     defer directory.Close()

     objects, err := directory.Readdir(-1)
     if err != nil {
     	return err
     }

     for _, obj := range objects {

         sourcefilepointer := filepath.Join(source, obj.Name())

         destinationfilepointer := filepath.Join(dest, obj.Name())


         if obj.IsDir() {
             // create sub-directories - recursively
             err = CopyDir(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println("CopyDir Error: ", err)
                 return err
             }
         } else {
             // perform copy
             err = CopyFile(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println("CopyFile Error: ", err)
                 return err
             }
         }

     }
     return nil
 }

  func CopyDir(source string, dest string) (err error) {
     sourceinfo, err := os.Stat(source)
     if err != nil {
         return err
     }

     err = os.MkdirAll(dest, sourceinfo.Mode())
     if err != nil {
         return err
     }

     directory, _ := os.Open(source)

     defer directory.Close()

     objects, err := directory.Readdir(-1)

     for _, obj := range objects {

         sourcefilepointer := filepath.Join(source, obj.Name())

         destinationfilepointer := filepath.Join(dest, obj.Name())


         if obj.IsDir() {
             // create sub-directories - recursively
             err = CopyDir(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println(err)
             }
         } else {
             // perform copy
             err = CopyFile(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println(err)
             }
         }

     }
     return
 }

