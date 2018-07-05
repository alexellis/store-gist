# store-gist

This function will store a gist using the GitHub Gist API. It uses OpenFaaS Cloud and an OAuth token sealed as a SealedSecret.

* Create the secret for your own account:

```
faas-cli cloud seal --name alexellis-store-gist --literal github-token=$TOKEN --cert=./pub-cert.pem
```

Example of posting to Gist adapted from blog post by [Minhazul Haque](https://bits.mdminhazulhaque.io/golang/create-gist-using-go.html)

This could also be achieved by vendoring the following library: https://github.com/google/go-github


