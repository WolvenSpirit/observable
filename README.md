[![observable](https://github.com/WolvenSpirit/observable/actions/workflows/go.yml/badge.svg)](https://github.com/WolvenSpirit/observable/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/WolvenSpirit/observable/branch/master/graph/badge.svg?token=TLWM6G8PWU)](https://codecov.io/gh/WolvenSpirit/observable)

### observable

#### Inspired by rxjs

- This package provides an `Observable` that can be shared between goroutines.
- The observable requires a pointer to a channel that is local to the specific goroutine and not shared, after the goroutine subscribes it can receive all notifications passed through the observable on the channel that it subscribed with.
- The package also has an `On()` hook method that accepts a value, signifying a possible event type, the goroutine local channel pointer and the callback that will be executed if that value gets sent through the observable.
