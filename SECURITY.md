# Security Policy

## Supported Versions

We currently support the following versions of ContainDB with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 4.12.x  | :white_check_mark: |
| 4.11.x  | :x:                |
| 4.10.x  | :x:                |
| < 4.0.0 | :x:                |

## Reporting a Vulnerability

ContainDB takes security issues seriously. We appreciate your efforts to responsibly disclose your findings and will make every effort to acknowledge your contributions.

### How to Report a Security Vulnerability

If you discover a security vulnerability within ContainDB, please send an email to ankansahaofficial@gmail.com with:

1. A detailed description of the vulnerability
2. Steps to reproduce the issue
3. The potential impact of the vulnerability
4. Any suggestions for remediation if available

Please **DO NOT** create public GitHub issues for security vulnerabilities.

### What to Expect

After reporting a vulnerability, you can expect:

- **Initial Response**: We'll acknowledge your email within 48 hours.
- **Validation Process**: We'll work to validate and reproduce the issue.
- **Resolution Timeline**: Once validated, we'll provide an estimated timeline for resolution.
- **Fix Implementation**: We'll implement a fix and test it thoroughly.
- **Public Disclosure**: After the fix is released, we may disclose the vulnerability with credit to you (unless you prefer to remain anonymous).

## Security Best Practices for ContainDB Users

1. **Always run with proper permissions**: While ContainDB needs sudo access to manage Docker containers, ensure you're using it in secure environments.

2. **Keep ContainDB updated**: Use the latest version to benefit from security fixes and improvements.

3. **Use secure database credentials**: When prompted for database credentials during setup, use strong, unique passwords.

4. **Be careful with network exposure**: Avoid exposing database ports to the internet unless absolutely necessary and properly secured.

5. **Regular backups**: While ContainDB helps manage your database containers, you should still maintain regular backups of important data.

## Container Security Considerations

ContainDB uses Docker containers to run databases. To enhance security:

- Consider running databases that need to be accessed externally behind a reverse proxy with TLS
- Review the Docker security documentation for additional hardening measures
- Regularly update the database images used by ContainDB

Thank you for helping keep ContainDB and its users safe!
