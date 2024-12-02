package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
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

func printBinaryData(file *os.File) {
	for {
		var eventType EventType
		err := binary.Read(file, binary.LittleEndian, &eventType)
		if err != nil {
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

			time.Sleep(1 * time.Second) //DELETE Убери если не нужны паузы на time.Stamp'ах.

		case EntityMove:
			var rplId int32
			var posX, posZ, rotY float32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &posX)
			binary.Read(file, binary.LittleEndian, &posZ)
			binary.Read(file, binary.LittleEndian, &rotY)
			fmt.Printf("EntityMove: RplId: %d, X: %.2f, Z: %.2f, RotY: %.2f\n", rplId, posX, posZ, rotY)

		case CharacterPossess:
			var rplId, playerId int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &playerId)
			fmt.Printf("CharacterPossess: RplId: %d, PlayerId: %d\n", rplId, playerId)

		case CharacterRegistration:
			var rplId, factionKeyLength int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &factionKeyLength)
			factionKey := make([]byte, factionKeyLength)
			file.Read(factionKey)
			fmt.Printf("CharacterRegistration: RplId: %d, FactionKey: %s\n", rplId, string(factionKey))

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
			fmt.Printf("VehicleRegistration: RplId: %d, VehicleName: %s, VehicleType: %d, FactionKey: %s\n",
				rplId, string(vehicleName), vehicleType, string(factionKey))

		case PlayerRegistration:

			var playerId, playerNameLen int32
			binary.Read(file, binary.LittleEndian, &playerId)
			binary.Read(file, binary.LittleEndian, &playerNameLen)
			playerName := make([]byte, playerNameLen)
			file.Read(playerName)
			fmt.Printf("PlayerRegistration: PlayerId: %d, PlayerName: %s\n", playerId, string(playerName))

		case EntityDamageStateChanged:
			var rplId, damageState int32
			binary.Read(file, binary.LittleEndian, &rplId)
			binary.Read(file, binary.LittleEndian, &damageState)
			fmt.Printf("EntityDamageStateChanged: RplId: %d, DamageState: %d\n", rplId, damageState)

		case CharacterBoardVehicle, CharacterUnBoardVehicle:
			var vehicleId, playerId int32
			binary.Read(file, binary.LittleEndian, &vehicleId)
			binary.Read(file, binary.LittleEndian, &playerId)
			action := "Board"
			if eventType == CharacterUnBoardVehicle {
				action = "UnBoard"
			}
			fmt.Printf("Character%sVehicle: VehicleId: %d, PlayerId: %d\n", action, vehicleId, playerId)

		case ProjectileShoot:
			var shootEntity int32
			var hitPosX, hitPosZ float32
			binary.Read(file, binary.LittleEndian, &shootEntity)
			binary.Read(file, binary.LittleEndian, &hitPosX)
			binary.Read(file, binary.LittleEndian, &hitPosZ)
			fmt.Printf("ProjectileShoot: Entity: %d, Hit X: %.2f, Hit Z: %.2f\n", shootEntity, hitPosX, hitPosZ)

		case Explosion:
			var hitPosX, hitPosZ, impulseDistance float32
			binary.Read(file, binary.LittleEndian, &hitPosX)
			binary.Read(file, binary.LittleEndian, &hitPosZ)
			binary.Read(file, binary.LittleEndian, &impulseDistance)
			fmt.Printf("Explosion: Hit X: %.2f, Hit Z: %.2f, Impulse Distance: %.2f\n", hitPosX, hitPosZ, impulseDistance)

		case EntityDelete:
			var rplId int32
			binary.Read(file, binary.LittleEndian, &rplId)
			fmt.Printf("EntityDelete: RplId: %d\n", rplId)

		default:
			fmt.Printf("Неизвестный тип события: %d\n", eventType)
			break
		}
	}
}
