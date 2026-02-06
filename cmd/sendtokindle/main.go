package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"sendtokindle/internal/httpapi"
	"sendtokindle/internal/storage"
	"sendtokindle/internal/web"
)

func main() {
	var port int
	var dir string
	flag.IntVar(&port, "port", 8080, "http server port")
	flag.StringVar(&dir, "dir", "", "storage directory (default: ~/.sendtokindle)")
	flag.Parse()

	var store *storage.Store
	var err error
	if dir != "" {
		store, err = storage.New(dir)
	} else {
		store, err = storage.NewDefault()
		if err != nil {
			fallback, ferr := storage.New("./.sendtokindle")
			if ferr != nil {
				log.Fatalf("init storage: %v", err)
			}
			log.Printf("sendtokindle: default storage unavailable: %v", err)
			store = fallback
			err = nil
		}
	}
	if err != nil {
		log.Fatalf("init storage: %v", err)
	}

	renderer, err := web.NewRenderer()
	if err != nil {
		log.Fatalf("init templates: %v", err)
	}

	kindleURL := bestKindleURL(port)

	handlers := &httpapi.Handlers{
		Store:     store,
		Renderer:  renderer,
		KindleURL: kindleURL,
		StoreRoot: store.Root(),
	}
	router := httpapi.NewRouter(httpapi.Config{}, handlers)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("sendtokindle: http://127.0.0.1:%d/  (admin: /admin)", port)
	log.Printf("sendtokindle: storage dir: %s", store.Root())
	if kindleURL != "" {
		log.Printf("sendtokindle: kindle url: %s", kindleURL)
	}

	go openBrowser(fmt.Sprintf("http://127.0.0.1:%d/admin", port))
	log.Fatal(srv.Serve(ln))
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

func localIPv4s() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	var ips []string
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP == nil {
			continue
		}
		ip := ipNet.IP.To4()
		if ip == nil || ip.IsLoopback() {
			continue
		}
		s := ip.String()
		if strings.HasPrefix(s, "169.254.") {
			continue
		}
		ips = append(ips, s)
	}
	return ips
}

func bestKindleURL(port int) string {
	candidates := bestIPv4Candidates()
	if len(candidates) == 0 {
		return ""
	}
	return fmt.Sprintf("http://%s:%d/", candidates[0], port)
}

func bestIPv4Candidates() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	type candidate struct {
		ip    string
		score int
	}
	var cands []candidate

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP == nil {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if isLinkLocalIPv4(ip) || isBenchmarkIPv4(ip) {
				continue
			}

			score := 0
			if isPrivateIPv4(ip) {
				score += 100
			}

			name := strings.ToLower(iface.Name)
			switch {
			case name == "en0":
				score += 50
			case name == "en1":
				score += 40
			case strings.HasPrefix(name, "wl"):
				score += 40
			}

			if strings.Contains(name, "utun") || strings.Contains(name, "vpn") {
				score -= 60
			}
			if strings.Contains(name, "bridge") || strings.Contains(name, "docker") || strings.Contains(name, "vbox") || strings.Contains(name, "vmnet") {
				score -= 40
			}

			s := ip.String()
			cands = append(cands, candidate{ip: s, score: score})
		}
	}

	if len(cands) == 0 {
		return nil
	}

	bestScore := cands[0].score
	for _, c := range cands[1:] {
		if c.score > bestScore {
			bestScore = c.score
		}
	}

	var best []string
	seen := make(map[string]struct{}, len(cands))
	for _, c := range cands {
		if c.score != bestScore {
			continue
		}
		if _, ok := seen[c.ip]; ok {
			continue
		}
		seen[c.ip] = struct{}{}
		best = append(best, c.ip)
	}
	if len(best) == 0 {
		return nil
	}
	return best
}

func isPrivateIPv4(ip net.IP) bool {
	if len(ip) != 4 {
		return false
	}
	switch {
	case ip[0] == 10:
		return true
	case ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31:
		return true
	case ip[0] == 192 && ip[1] == 168:
		return true
	default:
		return false
	}
}

func isLinkLocalIPv4(ip net.IP) bool {
	return len(ip) == 4 && ip[0] == 169 && ip[1] == 254
}

func isBenchmarkIPv4(ip net.IP) bool {
	if len(ip) != 4 {
		return false
	}
	if ip[0] == 198 && (ip[1] == 18 || ip[1] == 19) {
		return true
	}
	return false
}
