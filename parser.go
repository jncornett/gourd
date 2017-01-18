package main

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type Parser struct {
	s BufferedScanner
}

func NewParser(r io.Reader) Parser {
	s := bufio.NewScanner(r)
	fs := NewFilteredScanner(s, func(s string) bool {
		return len(trimString(s, "")) != 0
	})
	bs := NewBufferedScanner(fs)
	return Parser{s: bs}
}

func (p *Parser) Parse() ([]UseCase, error) {
	var useCases []UseCase
	for {
		useCase, err := p.parseUseCase()
		if err != nil {
			return nil, err
		}
		if useCase == nil {
			break
		}
		log.Println("Parsed a UseCase:", *useCase)
		useCases = append(useCases, *useCase)
	}
	return useCases, nil
}

func (p *Parser) parseUseCase() (*UseCase, error) {
	givens, err := p.parseBlocks("Given")
	if err != nil {
		return nil, err
	}
	if len(givens) == 0 {
		return nil, nil
	}
	log.Println("Parsed Some Givens:", givens)
	whenThens, err := p.parseWhenThens()
	if err != nil {
		return nil, err
	}
	log.Println("Parsed Some WhenThens:", whenThens)
	return &UseCase{Givens: givens, WhenThens: whenThens}, nil
}

func (p *Parser) parseWhenThens() ([]WhenThen, error) {
	var whenThens []WhenThen
	for {
		subWhenThens, err := p.parseWhenThen()
		if err != nil {
			return nil, err
		}
		if subWhenThens == nil {
			break
		}
		for _, whenThen := range subWhenThens {
			whenThens = append(whenThens, whenThen)
		}
		log.Println("Parsed some sub WhenThens:", subWhenThens)
	}
	return whenThens, nil
}

func (p *Parser) parseWhenThen() ([]WhenThen, error) {
	var whenThens []WhenThen
	when, err := p.parseBlock("When")
	if err != nil {
		return nil, err
	}
	if when == nil {
		return nil, nil
	}
	thens, err := p.parseBlocks("Then")
	if err != nil {
		return nil, err
	}
	for _, then := range thens {
		whenThens = append(
			whenThens,
			WhenThen{
				When: *when,
				Then: then,
			},
		)
	}
	return whenThens, nil
}

func (p *Parser) parseBlocks(prefix string) ([]Condition, error) {
	var blocks []Condition
	for {
		block, err := p.parseBlock(prefix)
		if err != nil {
			return nil, err
		}
		if block == nil {
			break
		}
		log.Println("Parsed a", prefix, "block:", block)
		blocks = append(blocks, *block)
	}
	return blocks, nil
}

func (p *Parser) parseBlock(prefix string) (*Condition, error) {
	for p.s.Scan() {
		line := p.s.Text()
		if strings.HasPrefix(line, prefix) {
			commands, err := p.parseCommands()
			if err != nil {
				return nil, err
			}
			cond := &Condition{
				Descr:    trimString(line, prefix),
				Commands: commands,
			}
			return cond, nil
		} else {
			p.s.Unscan()
			// return early instead of break so we don't confuse
			// with scanner error (which will also cause loop to exit)
			return nil, nil
		}
	}
	return nil, p.s.Err()
}

func (p *Parser) parseCommands() ([]string, error) {
	var (
		commands []string
		indent   int
	)
	for p.s.Scan() {
		line := p.s.Text()
		i := getIndent(line)
		if i >= indent {
			indent = i
			commands = append(commands, trimString(line, ""))
		} else {
			p.s.Unscan()
			// return early instead of break so we don't confuse
			// with scanner error (which will also cause loop to exit)
			return commands, nil
		}
	}
	err := p.s.Err()
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func trimString(s, prefix string) string {
	return strings.Trim(strings.TrimPrefix(s, prefix), " \t")
}

func getIndent(s string) int {
	for i, ch := range s {
		if ch != ' ' && ch != '\t' {
			return i
		}
	}
	return 0
}
