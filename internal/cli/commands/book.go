package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/ranking"
)

var filterOnCampus string

type roomFilter func(booking.Room) bool

func BookCmd(getBS func() booking.BookingService, getRS func() ranking.RankingService) *cobra.Command {
	bookCmd := &cobra.Command{
		Use:   "book {day} {time}",
		Short: "Create a booking",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, getBS, getRS)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if a := len(args); a > 2 || a == 1 {
				return fmt.Errorf("wrong number of arguments")
			}
			return nil
		},
	}
	bookCmd.Flags().StringVarP(&filterOnCampus, "campus", "c", "", "Show only rooms from either (J)ohanneberg or (L)indholmen")

	return bookCmd
}

func run(cmd *cobra.Command, args []string, getBS func() booking.BookingService, getRS func() ranking.RankingService) {
	bs := getBS()

	startDate, endDate := readArgs(args)

	available, err := bs.Available(startDate, endDate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rs := getRS()
	rankings, err := rs.GetRankings()
	if err != nil {
		fmt.Printf("Failed to get rankings: %v\n", err)
	} else {
		available = rankings.Sort(available)
	}

	var filters []roomFilter
	if filterOnCampus != "" {
		filters = append(filters, isOnCampus)
	}

	available = filterRooms(available, filters)

	showAvailable(available)

	room, err := prompt("Room to book")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	n, err := strconv.Atoi(room)
	n--
	if err != nil {
		fmt.Printf("Invalid Room\n")
		os.Exit(1)
	}

	message, err := prompt("Message to add with the booking (default: empty)")
	if err != nil {
		fmt.Println(err)
		fmt.Println("No booking was made")
		os.Exit(1)
	}

	if n < len(available) && n >= 0 {
		fmt.Printf("Booking %s...\n", available[n].Id)
		booking := booking.Booking{
			Room:  available[n],
			Start: startDate,
			End:   endDate,
			Text:  message,
		}
		err := bs.Book(booking)
		if err != nil {
			fmt.Println("couldn't book room")
			os.Exit(1)
		}
		fmt.Printf("Booked %s successfully!\n", available[n].Id)

		if rankings != nil {
			rankings.Update(available[n], available)
			err := rs.SaveRankings(rankings)
			if err != nil {
				fmt.Printf("Could not save updated rankings: %v\n", err)
			}
		}

	} else {
		print("no such booking")
	}
}

func filterRooms(rs []booking.Room, fs []roomFilter) []booking.Room {
	return filter(rs, func(r booking.Room) bool {
		for _, f := range fs {
			if !f(r) {
				return false
			}
		}
		return true
	})
}

func filter(xs []booking.Room, f roomFilter) []booking.Room {
	var ys []booking.Room
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

func isOnCampus(room booking.Room) bool {
	if filterOnCampus == "" {
		return true
	}
	if len(room.Campus) == 0 {
		return false
	}
	if len(filterOnCampus) == 1 {
		return string(strings.ToLower(room.Campus)[0]) == strings.ToLower(filterOnCampus)
	}
	return strings.ToLower(room.Campus) == strings.ToLower(filterOnCampus)
}

func prompt(message string) (string, error) {
	fmt.Printf("==> %s\n", message)
	fmt.Print("==> ")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("while getting input: %w", err)
	}
	return strings.Replace(input, "\n", "", -1), nil
}

func readArgs(args []string) (time.Time, time.Time) {
	date, err := extractDate(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	start, end, err := extractTimes(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return date.Add(start), date.Add(end)
}

func showAvailable(available []booking.Room) {
	for i := len(available) - 1; i >= 0; i-- {
		room := available[i]

		prevIsSame := i > 0 && available[i-1].Id == room.Id
		nextIsSame := i < len(available)-1 && available[i+1].Id == room.Id
		roomName := room.Id
		if prevIsSame || nextIsSame {
			roomName = fmt.Sprintf("%s.%s", room.Provider, room.Id)
		}

		fmt.Printf("%4s %-7s\n",
			fmt.Sprintf("[%d]", i+1),
			roomName,
		)
	}
}

func extractTimes(s string) (time.Duration, time.Duration, error) {
	parts := strings.Split(s, "-")
	if len(parts) == 2 {
		start, err := extractTime(parts[0])
		if err != nil {
			return time.Second, time.Second, err
		}
		end, err := extractTime(parts[1])
		if err != nil {
			return time.Second, time.Second, err
		}
		return start, end, nil
	} else {
		switch s {
		case "lunch":
			return time.Hour * 12, time.Hour * 13, nil
		default:
			return time.Second, time.Second, fmt.Errorf("failed to parse times from %s, s", s)
		}
	}
}

func extractTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	hour := ""
	minute := ""
	if len(parts) == 2 {
		hour = parts[0]
		minute = parts[1]
	} else if len(s) == 4 {
		hour = s[:2]
		minute = s[2:]
	} else if len(s) <= 2 {
		hour = s
		minute = "0"
	}
	h, err := strconv.Atoi(hour)
	if err != nil {
		return time.Nanosecond, err
	}
	m, err := strconv.Atoi(minute)
	if err != nil {
		return time.Nanosecond, err
	}
	return time.Hour*time.Duration(h) + time.Minute*time.Duration(m), nil
}

func extractDate(s string) (time.Time, error) {
	switch n := time.Now(); strings.ToLower(s) {
	case "today":
		return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location()), nil
	case "tomorrow":
		return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location()).Add(time.Hour * 24), nil
	default:
		t, err := extractDateAbsolute(s, n)
		if err != nil {
			return n, err
		}
		return t, nil
	}
}

func extractDateAbsolute(s string, n time.Time) (time.Time, error) {
	weekday, err := parseWeekday(strings.ToLower(s))
	if err == nil {
		diff := daysToAdd(n.Weekday(), weekday)
		t := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
		return t.Add(time.Hour * 24 * time.Duration(diff)), nil
	}

	switch len(s) {
	case 1:
		format := "2"
		t, err := time.ParseInLocation(format, s, n.Location())
		if err != nil {
			return n, err
		}
		t = t.AddDate(n.Year(), int(n.Month())-1, 0)
		if t.Day() < n.Day() {
			t = incMonth(t)
		}
		return t, nil
	case 2:
		format := "02"
		t, err := time.ParseInLocation(format, s, n.Location())
		if err != nil {
			return n, err
		}
		t = t.AddDate(n.Year(), int(n.Month())-1, 0)
		if t.Day() < n.Day() {
			t = incMonth(t)
		}
		return t, nil
	case 4:
		format := "0102"
		t, err := time.ParseInLocation(format, s, n.Location())
		if err != nil {
			return n, err
		}
		t = t.AddDate(n.Year(), 0, 0)
		if t.Month() < n.Month() || (t.Month() == n.Month() && t.Day() < n.Day()) {
			t = t.AddDate(1, 0, 0)
		}
		return t, nil
	case 6:
		format := "060102"
		t, err := time.ParseInLocation(format, s, n.Location())
		if err != nil {
			return n, err
		}
		return t, nil
	case 8:
		format := "20160102"
		t, err := time.ParseInLocation(format, s, n.Location())
		if err != nil {
			return n, err
		}
		return t, nil
	default:
		return n, fmt.Errorf("could not parse date from %s", s)
	}
}

func incMonth(t time.Time) time.Time {
	if t.Month() == 12 {
		return time.Date(t.Year()+1, 1, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	} else {
		return time.Date(t.Year(), t.Month()+1, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
}

func daysToAdd(from, to time.Weekday) int {
	d := len(daysOfWeek)
	daysToAdd := (int(to) - int(from) + d) % d
	if daysToAdd == 0 {
		daysToAdd += d
	}
	return daysToAdd
}

var daysOfWeek = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

func parseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}
	return time.Sunday, fmt.Errorf("invalid weekday format '%s'", v)
}
