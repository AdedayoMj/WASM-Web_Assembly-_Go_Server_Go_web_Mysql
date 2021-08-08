
The simplified and slightly modified version of [Lumberjack](https://github.com/natefinch/lumberjack) package that integrates nicely with the [Zap](https://github.com/uber-go/zap) sugarred logging which has a 4-10x faster than the other structured logging.

Both Zap and Lumberjack simplified by the [ITRLog](https://github.com/itrepablik/itrlog) package, so please use this ITRLog package instead.
 
# Installation
```
go get -u github.com/itrepablik/lumberjack
```

# Modified Version of Lumberjack
The main reason why we have some modified version of the Lumberjack package because of this line [#222](https://github.com/natefinch/lumberjack/blob/v2.0/lumberjack.go) from the source code of Lumberjack.  This function **openNew** will be executed during the log rotation which **renamed** the existing backup log file that exceeds the **MaxSize** in megabytes.

In our use case, we have constantly rotating backup scheduler program that keeps logging every time, this issue occurs as the file can't be renamed because it's open and use by another program error and that's a strange issue that even the existing close file function of Lumberjack can't close the existing used log file and then rename it.

At this point, we decided to clone the Lumberjack and modify it and replace the existing **os.Renamed** calls at this line #222: 

```
// Original os.Rename at line #222
if err := os.Rename(name, newname); err != nil {
	return fmt.Errorf("can't rename log file: %s", err)
}
```

and replace it with:
```
// Copy the current log file instead of just renaming it.
logDir := filepath.Dir(name)
dst := filepath.FromSlash(filepath.Join(logDir, filepath.Base(newname)))
if err := copyFile(name, dst, logDir); err != nil {
	return fmt.Errorf("can't backup the current log file: %s", err)
}
```

Besides, the **copyFile** function:
```
// copyFile copy a single file from the source to the destination.
func copyFile(src, dst, bareDst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		fmt.Println(err)
		return err
	}
	defer srcfd.Close()

	os.MkdirAll(bareDst, os.ModePerm) // Create the dst folder if not exist

	if dstfd, err = os.Create(dst); err != nil {
		fmt.Println(err)
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		fmt.Println(err)
		return err
	}

	if srcinfo, err = os.Stat(src); err != nil {
		fmt.Println(err)
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
```

# License
Code is distributed under MIT license, feel free to use it in your proprietary projects as well.
