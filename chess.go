package plugins

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"netsquirrel/utils"
)

type Chess struct {
	color  string 
	isTurn bool
	fen    string
	eng    *uci.Engine
}

func init() {
	Register("chess", &Chess{})
}

func (t *Chess) Description() string {
	return "Play a game of Chess against Stockfish."
}

func (t *Chess) Execute(comm Communicator, pluginDataChan chan<- string) {
    comm.Send("[Chess V1.0] Welcome to Chess! Type 'exit' to quit.")

    if rand.Intn(2) == 0 {
        t.color = "black"
    } else {
        t.color = "white"
    }
    comm.Send(fmt.Sprintf("[Chess V1.0] You are playing as %s.", t.color))

    eng, err := uci.New("stockfish")
    if err != nil {
        log.Fatalf("Failed to initialize Stockfish: %v", err)
    }
    defer eng.Close()
    t.eng = eng
    if err := t.eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
        log.Fatalf("Failed to start new game on Stockfish: %v", err)
    }

    t.fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    game := chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{}))
    t.isTurn = (t.color == "white")

    for {
        boardStr := t.DrawLargeBoard(game.Position().Board())
		comm.Send(boardStr)

        if t.isTurn {
            comm.Send("> ")
            input, err := comm.Receive()
            if err != nil {
                log.Printf("Error receiving input: %v", err)
                break
            }

            input = strings.TrimSpace(input)
            log.Printf("Received move input: %q", input)

            if input == "exit" {
                comm.Send("[Chess V1.0] Goodbye!")
                break
            }

            move, err := chess.AlgebraicNotation{}.Decode(game.Position(), input)
            if err != nil {
                comm.Send(fmt.Sprintf("[Chess V1.0] Invalid move: %s. Error: %v", input, err))
                continue
            }

            err = game.Move(move)
            if err != nil {
                comm.Send(fmt.Sprintf("[Chess V1.0] Failed to apply move: %s. Error: %v", input, err))
                continue
            }

            t.fen = game.FEN()
            t.isTurn = false
        } else {
            cmdPos := uci.CmdPosition{Position: game.Position()}
            cmdGo := uci.CmdGo{MoveTime: time.Second}
            if err := t.eng.Run(cmdPos, cmdGo); err != nil {
                log.Printf("Failed to get Stockfish move: %v", err)
                continue
            }

            move := t.eng.SearchResults().BestMove
            if err := game.Move(move); err != nil {
                log.Printf("Failed to make Stockfish move: %v", err)
                continue
            }

            t.fen = game.FEN()
            t.isTurn = true
            comm.Send(fmt.Sprintf("[Chess V1.0] Stockfish played: %s", move.String()))
        }

        if game.Outcome() != chess.NoOutcome {
            message := fmt.Sprintf("Game completed. %s by %s.", game.Outcome(), game.Method())
            comm.Send(message)
            break
        }
    }
}

func (t *Chess) DrawLargeBoard(board *chess.Board) string {
    pieceSymbols := map[chess.Piece]string{
        chess.NoPiece:        utils.ColorWithBackground("    ", utils.White, utils.Reset),
        chess.WhitePawn:      utils.ColorWithBackground(" ♙  ", utils.Black, utils.Pink),
        chess.WhiteKnight:    utils.ColorWithBackground(" ♘  ", utils.Black, utils.Pink),
        chess.WhiteBishop:    utils.ColorWithBackground(" ♗  ", utils.Black, utils.Pink),
        chess.WhiteRook:      utils.ColorWithBackground(" ♖  ", utils.Black, utils.Pink),
        chess.WhiteQueen:     utils.ColorWithBackground(" ♕  ", utils.Black, utils.Pink),
        chess.WhiteKing:      utils.ColorWithBackground(" ♔  ", utils.Black, utils.Pink),
        chess.BlackPawn:      utils.ColorWithBackground(" ♙  ", utils.White, utils.Orange),
        chess.BlackKnight:    utils.ColorWithBackground(" ♘  ", utils.White, utils.Orange),
        chess.BlackBishop:    utils.ColorWithBackground(" ♗  ", utils.White, utils.Orange),
        chess.BlackRook:      utils.ColorWithBackground(" ♖  ", utils.White, utils.Orange),
        chess.BlackQueen:     utils.ColorWithBackground(" ♕  ", utils.White, utils.Orange),
        chess.BlackKing:      utils.ColorWithBackground(" ♔  ", utils.White, utils.Orange),
    }

    bgColors := map[bool]string{
        true:  utils.BgBlack,
        false: utils.BgGrey,
    }

    var boardStr string
    boardStr += "   +--------------------------------+\n"

    rankStart, rankEnd, rankStep := 7, -1, -1
    fileOrder := []int{0, 1, 2, 3, 4, 5, 6, 7}
    if t.color == "black" {
        rankStart, rankEnd, rankStep = 0, 8, 1
        fileOrder = []int{7, 6, 5, 4, 3, 2, 1, 0}
    }

    for rank := rankStart; rank != rankEnd; rank += rankStep {
        boardStr += fmt.Sprintf(" %d |", rank+1)
        for _, file := range fileOrder {
            square := chess.Square(file + rank*8)
            piece := board.Piece(square)
            isDarkSquare := (file+rank)%2 == 0
            bgColor := bgColors[isDarkSquare]
    
            if piece == chess.NoPiece {
                boardStr += utils.ColorWithBackground("    ", "", bgColor)
            } else {
                symbol := pieceSymbols[piece]
                boardStr += utils.ColorWithBackground(symbol, "", bgColor)
            }
        }
        boardStr += "|\n"
    }

    boardStr += "   +--------------------------------+\n"
    if t.color == "white" {
        boardStr += "     a   b   c   d   e   f   g   h  \n"
    } else {
        boardStr += "     h   g   f   e   d   c   b   a  \n"
    }

    return boardStr
}








