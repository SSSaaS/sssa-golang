# sssaas-golang
An implementation of Shamir's Secret Sharing Algorithm in Go  

    Copyright (C) 2015 Alexander Scheel, Joel May, Matthew Burket  
    See Contributors.md for a complete list of contributors.  
    Licensed under the MIT License.  

## Usage
Note: this library is for a pure implementation of SSS in Go;
if you are looking for the API Library for SSSaaS, look [here](https://github.com/SSSAAS/sssaas-golang).

    sssa.create(minimum int, shares int, raw string) - creates a set of shares

    sssa.combine(shares []string) - combines shares into secret

For more detailed documentation, check out docs/sssa.md

## Contributing
We welcome pull requests, issues, security advice on this library, or other contributions you feel are necessary. Feel free to open an issue to discuss any questions you have about this library.

This is the reference implementation for all other SSSA projects.

For security issues, send a GPG-encrypted email to <alexander.m.scheel@gmail.com> with public key [0xBDC5F518A973035E](https://pgp.mit.edu/pks/lookup?op=vindex&search=0xBDC5F518A973035E).
