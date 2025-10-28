🎧 Spotify Featured Artists Dashboard

A full-stack web dashboard that highlights three random Spotify artists daily — helping users discover new music and emerging talent across 20+ genres.

The backend pulls live data from the Spotify Web API, generates a spotify_data.json file, and the frontend (HTML/CSS/JS) displays each artist’s name, image, genre, followers, top track, and Spotify link.
Updates are automated daily via GitHub Actions, and the site is deployed to GitHub Pages.

🖥️ Live Demo: https://waylonwang17.github.io/Spotify-Daily-Featured-Artists-Board

✨ Features

🎲 Randomly selects 3 genres and features one random artist from each every day

🎵 Displays artist name, followers, top track, genres, and Spotify link

🖼️ Clean visual layout with embedded Spotify logo and artist images

🔁 Daily automated refresh of artists using GitHub Actions (cron job)

🌐 Hosted on GitHub Pages — no server required

🔒 Uses Spotify’s Client Credentials flow (no user login or OAuth needed)

🧱 Tech Stack

Backend: Go
Frontend: HTML, CSS, JavaScript
API: Spotify Web API
Automation: GitHub Actions (CI/CD)
Hosting: GitHub Pages

🗺️ How It Works

The Go program requests an access token using the Spotify Client Credentials flow.

It picks 3 random genres and queries Spotify’s search API for matching artists.

For each artist, it retrieves:

Profile image

Genre(s)

Follower count

Top track name and link

It writes all this to a spotify_data.json file.

The frontend reads that JSON and displays the artist cards.

A GitHub Actions workflow runs daily to regenerate data and redeploy the site.
