window.ED or= {}

##
##
class window.ED.Card extends window.Malefic.Core
  ##
  ##
  constructor:() ->
    @bowed = false
    @_masked = true
    @tags = {}
    @traits = {}
    @PreLoad?()

  masked: ->
    @_masked

  ##
  ##
  unmask: (cardData) ->
    @_masked = false

  ######################
  ## START MATCH
  StageMatch_Phase(): ->
    @bowed = false
    if @masked() is false
      @Broker.Trigger('render:card', @)
  ######################
  ## START TURN
  Harvest_Phase(): ->

  ######################
  ##
  Event_Phase(): ->

  ######################
  ##
  Reaction_Phase(): ->

  ######################
  ##
  Conquest_Phase(): ->

  ######################
  ## END TURN
  ##
  Building_Phase(): ->

  ######################
  ##
  ## END MATCH
  PostMatch_Phase(): ->

##
##
##

class window.ED.GodCard extends window.ED.Card
  PreLoad: ->

class window.ED.FortressCard extends window.ED.Card
  PreLoad: ->

# class window.ED.LeaderCard extends window.ED.Card
#   PreLoad: ->

##
##
##

class window.ED.DestinyCard extends window.ED.Card
  PreLoad: ->

class window.ED.EmpireCard extends window.ED.Card
  PreLoad: ->
