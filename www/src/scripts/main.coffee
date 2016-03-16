window.Malefic._broker = new window.Malefic.Stream()

class window.Core extends window.Malefic.Core
  constructor:  ->
    @widget = {}
    @Broker = window.Malefic._broker

    @widget.login = new window.ED.PlayingTable()
    @widget.login.Ready( ->
      console.log('PlayingTable Loaded')
    )


core = new window.Core()
console.log(core)
