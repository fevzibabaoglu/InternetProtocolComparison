---

# InternetProtocolComparison
Comparison of TCP/UDP/QUIC by speed and reliability in GO.

---

## The results gathered on my end are the following:

|                   | TCP - (Bidirectional) | UDP - (Unidirectional) | QUIC by [quic-go](https://pkg.go.dev/github.com/quic-go/quic-go) - (Bidirectional) |
| ----------------- | --------------------- | ---------------------- | ---------------------------------------------------------------------------------- |
| **Total***        | 9.6897449s            | 13.2409361s            | 4.0459975s                                                                         |
| **Time per byte** | 37ns                  | 51ns                   | 15ns                                                                               |

**Total***: 1000000 messages with a length of 255 bytes (total of 255000000 bytes) 

---

## Notes:
The timings were calculated only in localhost. The results will vary depending on usage of the internet, wifi, etc.

---
