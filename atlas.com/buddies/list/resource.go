package list

import (
	"atlas-buddies/buddy"
	list2 "atlas-buddies/kafka/message/list"
	"atlas-buddies/kafka/producer"
	list3 "atlas-buddies/kafka/producer/list"
	"atlas-buddies/rest"
	"errors"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-rest/server"
	"github.com/gorilla/mux"
	"github.com/jtumidanski/api2go/jsonapi"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

const (
	GetBuddyList          = "get_buddy_list"
	CreateBuddyList       = "create_buddy_list"
	GetBuddiesInBuddyList = "get_buddies_in_buddy_list"
	AddBuddyToBuddyList   = "add_buddy_to_buddy_list"
	UpdateBuddyListCapacity = "update_buddy_list_capacity"
)

func InitResource(si jsonapi.ServerInformation) func(db *gorm.DB) server.RouteInitializer {
	return func(db *gorm.DB) server.RouteInitializer {
		return func(router *mux.Router, l logrus.FieldLogger) {
			registerGet := rest.RegisterHandler(l)(si)
			r := router.PathPrefix("/characters/{characterId}/buddy-list").Subrouter()
			r.HandleFunc("", registerGet(GetBuddyList, handleGetBuddyList(db))).Methods(http.MethodGet)
			r.HandleFunc("", rest.RegisterInputHandler[RestModel](l)(si)(CreateBuddyList, handleCreateBuddyList)).Methods(http.MethodPost)
			r.HandleFunc("/buddies", registerGet(GetBuddiesInBuddyList, handleGetBuddiesInBuddyList(db))).Methods(http.MethodGet)
			r.HandleFunc("/buddies", rest.RegisterInputHandler[buddy.RestModel](l)(si)(AddBuddyToBuddyList, handleAddBuddyToBuddyList)).Methods(http.MethodPost)
			r.HandleFunc("/capacity", rest.RegisterInputHandler[RestModel](l)(si)(UpdateBuddyListCapacity, handleUpdateBuddyListCapacity)).Methods(http.MethodPut)
		}
	}
}

func handleGetBuddyList(db *gorm.DB) rest.GetHandler {
	return func(d *rest.HandlerDependency, c *rest.HandlerContext) http.HandlerFunc {
		return rest.ParseCharacterId(d.Logger(), func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				bl, err := NewProcessor(d.Logger(), d.Context(), db).GetByCharacterId(characterId)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				res, err := model.Map(Transform)(model.FixedProvider(bl))()
				if err != nil {
					d.Logger().WithError(err).Errorf("Creating REST model.")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				server.Marshal[RestModel](d.Logger())(w)(c.ServerInformation())(res)
			}
		})
	}
}

func handleCreateBuddyList(d *rest.HandlerDependency, _ *rest.HandlerContext, i RestModel) http.HandlerFunc {
	return rest.ParseCharacterId(d.Logger(), func(characterId uint32) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := producer.ProviderImpl(d.Logger())(d.Context())(list2.EnvCommandTopic)(list3.CreateCommandProvider(characterId, i.Capacity))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusAccepted)
		}
	})
}

func handleGetBuddiesInBuddyList(db *gorm.DB) rest.GetHandler {
	return func(d *rest.HandlerDependency, c *rest.HandlerContext) http.HandlerFunc {
		return rest.ParseCharacterId(d.Logger(), func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				bl, err := NewProcessor(d.Logger(), d.Context(), db).GetByCharacterId(characterId)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				res, err := model.SliceMap(buddy.Transform)(model.FixedProvider(bl.buddies))()()
				if err != nil {
					d.Logger().WithError(err).Errorf("Creating REST model.")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				server.Marshal[[]buddy.RestModel](d.Logger())(w)(c.ServerInformation())(res)
			}
		})
	}
}

func handleAddBuddyToBuddyList(d *rest.HandlerDependency, _ *rest.HandlerContext, i buddy.RestModel) http.HandlerFunc {
	return rest.ParseCharacterId(d.Logger(), func(characterId uint32) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//err := producer.ProviderImpl(d.Logger())(d.Context())(EnvCommandTopic)(addBuddyCommandProvider(characterId, i.CharacterId, i.Group, i.CharacterName, i.ChannelId, i.Visible))
			//if err != nil {
			//	w.WriteHeader(http.StatusInternalServerError)
			//	return
			//}

			w.WriteHeader(http.StatusAccepted)
		}
	})
}

func handleUpdateBuddyListCapacity(d *rest.HandlerDependency, _ *rest.HandlerContext, i RestModel) http.HandlerFunc {
	return rest.ParseCharacterId(d.Logger(), func(characterId uint32) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Validate capacity
			if i.Capacity == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Use default worldId for REST endpoints
			worldId := byte(1)

			// Send capacity update command
			err := producer.ProviderImpl(d.Logger())(d.Context())(list2.EnvCommandTopic)(UpdateCapacityCommandProvider(characterId, worldId, i.Capacity))
			if err != nil {
				d.Logger().WithError(err).Errorf("Error sending capacity update command for character [%d]", characterId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusAccepted)
		}
	})
}
