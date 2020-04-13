# Cryptographic contact tracing implementation

Implementation of the Apple and Google contact tracing cryptographic [specification](https://covid19-static.cdn-apple.com/applications/covid19/current/static/contact-tracing/pdf/ContactTracing-CryptographySpecification.pdf).

## Why?
I wondered how well the specification would scale, SHA-256 (and therefore HMAC-SHA265) is relatively slow, so given a large amount of daily keys it could take a lot of processing power to generate all the proximity keys to check for matches.

I wanted to test how long it would take to generate proximity keys for a list of daily keys.

### Findings
> **Please note:**
> 
> ⚠️ The number of infected people and the number of daily keys uploaded per person were not chosen using scientific methods.
>
> ⚠️ This implementation is not the most efficient.
>
> ⚠️ The implementation may not be correct, I believe I've followed the specification, but I can't find any examples to test against.

I assumed that 81,000 new people would be infected and for each infected person the daily keys for the past 14 days would be published.

I ran the test on my 2019 Macbook Pro (Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz).

When running on all CPUs it takes around 28 seconds to generate all of the proximity keys.

When running on a single processor it took around 120 seconds.
