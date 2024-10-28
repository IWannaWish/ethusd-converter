# /pkg
The /pkg directory contains library code intended for use by external applications (e.g., /pkg/mypubliclib). Other projects will import these libraries with the expectation that they function correctly, so be sure to carefully consider what you include here.

Keep in mind that using the internal directory is a more effective way to prevent other packages from importing your private code, as this restriction is enforced by Go. However, the /pkg directory is still a clear way to indicate that the code within is safe for external use. For a deeper understanding of the differences between the pkg and internal directories and guidance on when to use them, refer to Travis Jeffery's blog post, "I'll take pkg over internal."
