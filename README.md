---

# InternetProtocolComparison
Comparison of TCP/UDP/QUIC by speed and reliability in GO.

---

## The results gathered on my end are the following:

### [TCP] - (Bidirectional)
1000000 messages with a length of 255 bytes (total of 255000000 bytes) interchanged in 9.6897449s (byte per 37ns).

### [UDP] - (Unidirectional)
1000000 messages with a length of 255 bytes (total of 255000000 bytes) interchanged in 13.2409361s (byte per 51ns).

### [QUIC by [quic-go](https://pkg.go.dev/github.com/quic-go/quic-go)] - (Bidirectional)
1000000 messages with a length of 255 bytes (total of 255000000 bytes) interchanged in 4.0459975s (byte per 15ns).

---

## Notes:
The timings were calculated only in localhost. The results will vary depending on usage of the internet, wifi, etc.

---
