GopherCon Lightning Talk 2017 example code
==

This is the example API I presented at GopherCon 2017 during the lightning
talks. The API was written in about an hour for this talk. It's has a few
caveats:

* The API is not stable, but can be used as a starting point.
* The Responser serves multiple purposes, modifying the ResposeWriter
  & transforming the data. Ideally, these would be done separately.
* Any transforms to the data happen in order, so order matters. So if you're
  returning transformed data to be written, that should be last in the
  transform pipeline.
