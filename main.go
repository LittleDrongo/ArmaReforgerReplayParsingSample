package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LittleDrongo/fmn-lib/utils/jsn"
)

type EventType uint8

const (
	WorldTime EventType = iota
	EntityMove
	CharacterPossess
	CharacterRegistration
	VehicleRegistration
	PlayerRegistration
	EntityDamageStateChanged
	CharacterBoardVehicle
	CharacterUnBoardVehicle
	ProjectileShoot
	Explosion
	EntityDelete
)

func main() {
	file, err := os.Open("./replays/Replay1732461188.bin")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	printBinaryData(file)
}

type DataStorage struct {
	Events []Event[any] `json:"events"`
}
type Event[T any] struct {
	EventType EventType `json:"event_type"`
	EventData T         `json:"event_data"`
}

var storage DataStorage

func printBinaryData(file *os.File) {
	for {
		var eventType EventType
		err := binary.Read(file, binary.LittleEndian, &eventType)
		if err != nil {
			err = jsn.Export(storage, "./data/out/replay.json")
			if err != nil {
				log.Println("Шеф всё пропало!", err)
			}
			break // Конец файла
		}

		switch eventType {
		case WorldTime:
			var timeMS int32
			binary.Read(file, binary.LittleEndian, &timeMS)

			// Преобразуем миллисекунды в секунды и получаем значение времени
			duration := time.Duration(timeMS) * time.Millisecond
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			seconds := int(duration.Seconds()) % 60
			timeHHMMSS := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

			fmt.Printf("WorldTime: %d ms %s \n", timeMS, timeHHMMSS)

			type tempData struct {
				TimeMS int32 `json:"time_ms"`
			}

			event := Event[any]{
				EventType: WorldTime,
				EventData: tempData{
					TimeMS: timeMS,
				},
			}

			storage.Events = append(storage.Events, event)

			// time.Sleep(100 * time.Microsecond) //DELETE Убери если не нужны паузы на time.Stamp'ах.

		case EntityMove:
			var rplId int32
			var posX, posZ, rotY float32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &posX)
			binary.Read(file, binary.LittleEndian, &posZ)
			binary.Read(file, binary.LittleEndian, &rotY)

			type tempData struct {
				PosX  float32 `json:"pos_x"`
				PosZ  float32 `json:"pos_y"`
				RotY  float32 `json:"rot_y"`
				RplId int32   `json:"rpl_id"`
			}

			event := Event[any]{
				EventType: EntityMove,
				EventData: tempData{
					RplId: rplId,
					PosX:  posX,
					PosZ:  posZ,
					RotY:  rotY,
				},
			}

			storage.Events = append(storage.Events, event)

			// fmt.Printf("EntityMove: RplId: %d, X: %.2f, Z: %.2f, RotY: %.2f\n", rplId, posX, posZ, rotY)

		case CharacterPossess:
			var rplId, playerId int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &playerId)

			{
				type tempData struct {
					RplId    int32 `json:"rpl_id"`
					PlayerId int32 `json:"player_id"`
				}

				event := Event[any]{
					EventType: CharacterPossess,
					EventData: tempData{
						RplId:    rplId,
						PlayerId: playerId,
					},
				}
				storage.Events = append(storage.Events, event)
			}

			// fmt.Printf("CharacterPossess: RplId: %d, PlayerId: %d\n", rplId, playerId)

		case CharacterRegistration:
			var rplId, factionKeyLength int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &factionKeyLength)
			factionKey := make([]byte, factionKeyLength)
			file.Read(factionKey)

			{
				type tempData struct {
					FactionKeyLength int32
					RplId            int32
				}

				event := Event[any]{
					EventType: CharacterRegistration,
					EventData: tempData{
						RplId:            rplId,
						FactionKeyLength: factionKeyLength,
					},
				}
				storage.Events = append(storage.Events, event)
			}

			// fmt.Printf("CharacterRegistration: RplId: %d, FactionKey: %s\n", rplId, string(factionKey))

		case VehicleRegistration:
			var rplId, vehicleNameLen, vehicleType, factionLen int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &vehicleNameLen)
			vehicleName := make([]byte, vehicleNameLen)
			file.Read(vehicleName)
			binary.Read(file, binary.LittleEndian, &vehicleType)
			binary.Read(file, binary.LittleEndian, &factionLen)
			factionKey := make([]byte, factionLen)
			file.Read(factionKey)
			// fmt.Printf("VehicleRegistration: RplId: %d, VehicleName: %s, VehicleType: %d, FactionKey: %s\n",
			// 	rplId, string(vehicleName), vehicleType, string(factionKey))

			{
				type tempData struct {
					FactionKeyLength int32
					RplId            int32
					VehicleNameLen   int32
					VehicleType      int32
					FactionLen       int32
				}

				event := Event[any]{
					EventType: VehicleRegistration,
					EventData: tempData{
						RplId:          rplId,
						VehicleNameLen: vehicleNameLen,
						VehicleType:    vehicleType,
						FactionLen:     vehicleNameLen,
					},
				}
				storage.Events = append(storage.Events, event)
			}

		case PlayerRegistration:

			var playerId, playerNameLen int32
			binary.Read(file, binary.LittleEndian, &playerId)
			binary.Read(file, binary.LittleEndian, &playerNameLen)
			playerName := make([]byte, playerNameLen)
			file.Read(playerName)

			{
				type tempData struct {
					PlayerId      int32
					PlayerNameLen int32
				}

				event := Event[any]{
					EventType: PlayerRegistration,
					EventData: tempData{
						PlayerId:      playerId,
						PlayerNameLen: playerNameLen,
					},
				}
				storage.Events = append(storage.Events, event)
			}
			// fmt.Printf("PlayerRegistration: PlayerId: %d, PlayerName: %s\n", playerId, string(playerName))

		case EntityDamageStateChanged:
			var rplId, damageState int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &damageState)

			{
				type tempData struct {
					RplID       int32
					DamageState int32
				}

				event := Event[any]{
					EventType: EntityDamageStateChanged,
					EventData: tempData{
						RplID:       rplId,
						DamageState: damageState,
					},
				}
				storage.Events = append(storage.Events, event)
			}

			// fmt.Printf("EntityDamageStateChanged: RplId: %d, DamageState: %d\n", rplId, damageState)

		case CharacterBoardVehicle, CharacterUnBoardVehicle:
			var vehicleId, playerId int32
			binary.Read(file, binary.LittleEndian, &vehicleId)
			binary.Read(file, binary.LittleEndian, &playerId)
			// action := "Board"
			// if eventType == CharacterUnBoardVehicle {
			// 	action = "UnBoard"
			// }

			{

				type action int32
				const (
					actionBoard = iota
					actionUnboard
				)

				type tempData struct {
					VehicleId int32
					PlayerId  int32
					Action    action
				}

				var act action = actionUnboard
				var evType EventType = CharacterBoardVehicle
				if eventType == CharacterBoardVehicle {
					act = actionUnboard
					evType = CharacterUnBoardVehicle
				}

				event := Event[any]{
					EventType: evType,
					EventData: tempData{
						VehicleId: vehicleId,
						PlayerId:  playerId,
						Action:    act,
					},
				}
				storage.Events = append(storage.Events, event)
			}
			// fmt.Printf("Character%sVehicle: VehicleId: %d, PlayerId: %d\n", action, vehicleId, playerId)

		case ProjectileShoot:
			var shootEntity int32
			var hitPosX, hitPosZ float32
			binary.Read(file, binary.LittleEndian, &shootEntity)
			binary.Read(file, binary.LittleEndian, &hitPosX)
			binary.Read(file, binary.LittleEndian, &hitPosZ)

			{
				type tempData struct {
					ShootEntity int32
					HitPosX     float32
					HitPosZ     float32
				}

				event := Event[any]{
					EventType: ProjectileShoot,
					EventData: tempData{
						ShootEntity: shootEntity,
						HitPosX:     hitPosX,
						HitPosZ:     hitPosZ,
					},
				}
				storage.Events = append(storage.Events, event)
			}

			// fmt.Printf("ProjectileShoot: Entity: %d, Hit X: %.2f, Hit Z: %.2f\n", shootEntity, hitPosX, hitPosZ)

		case Explosion:
			var hitPosX, hitPosZ, impulseDistance float32
			binary.Read(file, binary.LittleEndian, &hitPosX)
			binary.Read(file, binary.LittleEndian, &hitPosZ)
			binary.Read(file, binary.LittleEndian, &impulseDistance)

			{
				type tempData struct {
					ImpulseDistance float32
					HitPosX         float32
					HitPosZ         float32
				}

				event := Event[any]{
					EventType: Explosion,
					EventData: tempData{
						ImpulseDistance: impulseDistance,
						HitPosX:         hitPosX,
						HitPosZ:         hitPosZ,
					},
				}
				storage.Events = append(storage.Events, event)
			}

			// fmt.Printf("Explosion: Hit X: %.2f, Hit Z: %.2f, Impulse Distance: %.2f\n", hitPosX, hitPosZ, impulseDistance)

		case EntityDelete:
			var rplId int32
			binary.Read(file, binary.LittleEndian, &rplId)

			{
				type tempData struct {
					RplId int32
				}

				event := Event[any]{
					EventType: EntityDelete,
					EventData: tempData{
						RplId: rplId,
					},
				}
				storage.Events = append(storage.Events, event)
			}
			// fmt.Printf("EntityDelete: RplId: %d\n", rplId)

		default:
			fmt.Printf("Неизвестный тип события: %d\n", eventType)
			break
		}
	}
}
