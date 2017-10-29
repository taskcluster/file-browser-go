package browser

// Operations we support (or plan to support in the future)
// READ - Similar to read()
// WRITE - Similar to write()
// TRUNC - Truncate a file to a given size.
// CREATE - Create a new file.
// REMOVE - Remove a file or directory.
// RENAME - Rename a file or directory.
// MKDIR - Create a directory.
// READDIR - This actually works similar to `ls`. You still get dirents (or something close enough anyway).
// STAT - Call stat(). If not supported by OS return `os.FileInfo` wrapped in an `attr` object.
//
// Extended operations:
// CP: Copy could have been implemented using bindings, but the data would require to take a much longer path.
// [ fs -> client -> fs ] CP should not depend on client connection. Therefore CP is implemented as an operation.
//
// What we don't support:
// OPEN: The browser is stateless and doesn't store any file descriptors. Therefore, just opening a file is
// meaningless. Use READ, WRITE, TRUNC, or CREATE directly.
// INOTIFY: We may support this in a future version. This would contain more than 1 command and will make the
// browser stateful.
// SETXATTR/GETXATTR: These are not supported.
// MKNODE: Create a file or a directory. Creating inodes is not supported.
// MOUNT/UNMOUNT
// LINK
// STAT_FS
//
// These lists are subject to change.
