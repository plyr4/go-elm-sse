port module Main exposing (..)

import Browser
import Html exposing (Html, button, div, h1,p, img,code, text)
import Html.Attributes exposing (src)
import Html.Events exposing (onClick)



---- MODEL ----


type alias Model =
    { data : Maybe String }

type alias  SSEConfig = {
        url : String
    } 
sseEndpoint : String
sseEndpoint =
    -- "http://localhost:8080/DavidVader/heyvela/builds/2/steps/3/logs/events"
    "http://localhost:9999/stream"

init : ( Model, Cmd Msg )
init =
    ( { data = Nothing }, Cmd.none )



---- UPDATE ----


type Msg
    = ConnectSSE
    | ReceiveEvent String



update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        ConnectSSE ->
            ( model, connectSSE <| SSEConfig sseEndpoint )

        ReceiveEvent message ->
            ( { model | data = Just (Maybe.withDefault "" model.data ++ "\n" ++ message) }, Cmd.none )



---- VIEW ----


view : Model -> Html Msg
view model =
    let

        data = Maybe.withDefault "" model.data
    in
    div []
        [ h1 [] [ text "sse proof of concept"]
        , div [] [ button [ onClick ConnectSSE ] [ text "connect" ] ]
        , p [][ code [] [text data]]

        ]



---- PROGRAM ----


main : Program () Model Msg
main =
    Browser.element
        { view = view
        , init = \_ -> init
        , update = update
        , subscriptions = subscriptions
        }


subscriptions : Model -> Sub Msg
subscriptions _ =
    messageReceiver ReceiveEvent


port connectSSE : SSEConfig -> Cmd msg


port messageReceiver : (String -> msg) -> Sub msg
