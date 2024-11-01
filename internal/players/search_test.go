package players_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchByCountry(t *testing.T) {
	// Given
	ctx := context.TODO()
	searchCriteria := players.SearchCriteria{
		Country: players.NewString("UK"),
	}

	wantPlayerList := unittests.SearchPlayersResultFixture(t)

	searchResult := players.SearchResult{
		Total:  27,
		Items:  wantPlayerList,
		Limit:  3,
		Offset: 0,
	}

	want := players.SearchResult{
		Total:  27,
		Items:  wantPlayerList,
		Limit:  3,
		Offset: 0,
	}

	storageMock := unittests.NewStorageMock()
	storageMock.On("Search", ctx, mock.AnythingOfType("players.SearchCriteria")).Return(&searchResult, nil)

	service, _ := unittests.NewPlayerServiceWithStorage(storageMock)

	// When
	got, err := service.List(ctx, searchCriteria)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, &want, got)
}

func TestSearchByCountryButError(t *testing.T) {
	// Given
	ctx := context.TODO()
	searchCriteria := players.SearchCriteria{
		Country: players.NewString("UK"),
	}

	want := "unable to list players: unable to search players: unexpected search error"

	searchError := errors.New("unable to search players: unexpected search error")

	storageMock := unittests.NewStorageMock()
	storageMock.On("Search", ctx, mock.AnythingOfType("players.SearchCriteria")).Return(nil, searchError)

	service, _ := unittests.NewPlayerServiceWithStorage(storageMock)

	// When
	got, err := service.List(ctx, searchCriteria)

	// Then
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Equal(t, want, err.Error())
}

func TestSearchPlayersWithEndpointSuccessfully(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	searchCriteria := players.SearchCriteria{
		Country: players.NewString("UK"),
	}

	wantPlayerList := unittests.SearchPlayersResultFixture(t)

	searchResult := players.SearchResult{
		Total:  27,
		Items:  wantPlayerList,
		Limit:  3,
		Offset: 0,
	}

	want := players.SearchPlayersDataResult{
		SearchResult: &players.SearchResult{
			Total:  27,
			Items:  wantPlayerList,
			Limit:  3,
			Offset: 0,
		},
	}

	storageMock := unittests.NewStorageMock()
	storageMock.On("Search", ctx, mock.AnythingOfType("players.SearchCriteria")).Return(&searchResult, nil)

	service, logger := unittests.NewPlayerServiceWithStorage(storageMock)
	searchPlayersEndpoint := players.MakeSearchPlayersEndpoint(service, logger)

	// When
	got, err := searchPlayersEndpoint.Do(ctx, searchCriteria)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)
}

func TestSearchPlayersWithEndpointButError(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	searchCriteria := players.SearchCriteria{
		Country: players.NewString("UK"),
	}

	want := players.SearchPlayersDataResult{
		Err: "unable to list players: unable to search players: unexpected search error",
	}

	searchError := errors.New("unable to search players: unexpected search error")

	storageMock := unittests.NewStorageMock()
	storageMock.On("Search", ctx, mock.AnythingOfType("players.SearchCriteria")).Return(nil, searchError)

	service, logger := unittests.NewPlayerServiceWithStorage(storageMock)
	searchPlayersEndpoint := players.MakeSearchPlayersEndpoint(service, logger)

	// When
	got, err := searchPlayersEndpoint.Do(ctx, searchCriteria)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
