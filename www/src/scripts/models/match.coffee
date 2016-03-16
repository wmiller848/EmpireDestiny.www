class window.ED.Match extends window.Malefic.Core

  constructor: ->


  StageMove: (move)->
    

  CommitMoves: ->
    url = 'https://api.ed.com/match/'
    request = @Ajax(url, "POST", data, headers)
    request.Progress = (status) ->
      console.log(status) # '0%', '2%', '10%', etc..
    request.Then = (err, response) ->
      console.log(err, response)
      # Attempt to get parsed JSON
      json = response.toJSON()
