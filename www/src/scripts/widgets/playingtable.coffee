window.ED or= {}

class window.ED.PlayingTable extends window.Malefic.View

  Context: '[data-id="playingtable_container"]'

  Template: 'playingtable.html.tmpl'

  Events:
    'hide:playingtable': 'hide'
    'show:playingtable': 'show'
    'lock:playingtable': 'lock'
    'unlock:playingtable': 'unlock'

  Data:
    'Title': 'playingtable',
    'Model': null

  Helpers:
    'log': ->
      @Log(arguments)

  Actions: ->
    'default': =>
      @Log(@)
    'fullscreen': =>
      @ToggleFullScreen()
    'hide': =>
      @Hide()
    'show': =>
      @Show()

  Elements:
    'playingtable': '[data-id="playingtable"]'

  Loaded: ->
    @Log('PlayingTable Widget Loaded')
    #@Hide()

    # Make an instance of two and place it on the page.
    el = @Q('[data-id="playingtable_container"]')
    params = { fullscreen: true }
    console.log(params)
    @two = new Two(params).appendTo(el)

    # circle = two.makeCircle(-70, 0, 50)
    # rect = two.makeRectangle(70, 0, 100, 100)
    # circle.fill = '#FF8000'
    # rect.fill = 'rgba(0, 200, 255, 0.75)'
    #
    # group = two.makeGroup(circle, rect)
    # group.translation.set(two.width / 2, two.height / 2)
    # group.scale = 0
    # group.noStroke()
    @two.bind('update', @Play).play() # Finally, start the animation loop

    # @Broker.On('event', (collections) =>
    #   @Render()
    # )

  OnBind: ->
    @Log('PlayingTable Binded')

    @Elements.playingtable?.on('click', (e) =>
      return if e.isTriggered
      e.preventDefault()
      e.stopPropagation()
      e.isTriggered = true
    )

  Play: (frameCount) ->
    # if group.scale > 0.9999
    #   group.scale = group.rotation = 0
    # t = (1 - group.scale) * 0.125
    # group.scale += t
    # group.rotation += t * 4 * Math.PI
