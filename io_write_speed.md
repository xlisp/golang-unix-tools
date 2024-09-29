## ubuntu test io write speed 

You can test the disk I/O write speed in Ubuntu using a few different methods. One of the most common tools is `dd`. Here's how you can do it:

### Test with `dd`:

Open the terminal and run the following command:

```bash
dd if=/dev/zero of=tempfile bs=1M count=1024 conv=fdatasync,notrunc
```

### Explanation:
- `if=/dev/zero`: This uses `/dev/zero` as the input file, which is a special file that provides null bytes (zeros).
- `of=tempfile`: This specifies the output file as `tempfile`. This is where the data is written to. You can delete this file after the test.
- `bs=1M`: This sets the block size to 1MB.
- `count=1024`: This writes 1024 blocks of data (so it writes 1GB of data).
- `conv=fdatasync,notrunc`: This ensures that the data is actually written to the disk, rather than just being cached by the OS.

### Output Example:
```bash
1024+0 records in
1024+0 records out
1073741824 bytes (1.1 GB, 1.0 GiB) copied, 3.53233 s, 304 MB/s
```

In this example, the write speed is 304 MB/s.

### Clean up:
After the test, you can remove the temporary file with:

```bash
rm tempfile
```

### Alternative Method:
You can also use a tool like `fio` for more advanced I/O benchmarking. Install `fio` using:

```bash
sudo apt install fio
```

Then run a test like this:

```bash
fio --name=write_test --size=1G --filesize=1G --filename=testfile --bs=1M --nrfiles=1 --direct=1 --sync=1 --rw=write --iodepth=1 --numjobs=1
```

This will give you more detailed output on write speeds and latencies.

