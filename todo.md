In this exercise we are taking a web application inspired by Tom Speak's _Quiet Hacker News_ and working on adding concurrency to it.

The initial version of the application each story is retrieved sequentially, meaning that before we can retrieve the 7th story, we must first retrieve the 1st-6th stories.

Our goal is to update the application to permit us to retrieve an arbitrary number of stories concurrently given the following limitations:

1. The order of the stories must be retained in the final output.
2. Discussions must be excluded from the final output (determined via the `isDiscussion` method on the `story` type).
3. We must return exactly `numStories` stories, assuming there are enough valid stories to do so.

What makes this problem challenging is that we can't know exactly how many stories we will need to retrieve upfront given the exclusion of discussions, and we also need to retain the order which means we can't just use the first `numStories` stories retrieved. It is possible that stories 2-31 will be retrieved prior to story 1.

The primary goals in this excercise are:

1. To practice working with concurrency (goroutines, channels, etc) in Go
2. To demonstrate that adding concurrency to an application isn't simply adding these things, but can often involve rethinking a problem entirely.

Note: _Quiet Hacker News_ is a project that Tom Speak is working on while learning Go. You can read more about it here: https://tomspeak.co.uk/posts/quiet-hacker-news
