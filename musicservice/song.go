package musicservice

type Song struct {
	Id        int      `json:"id"`
	Videoid   string   `json:"videoid"`
	Name      string   `json:"name"`
	Length    int      `json:"length"`
	Seek      int      `json:"seek"`
	AddedBy   string   `json:"added_by"`
	Thumbnail string   `json:"thumbnail"`
	Details   SongInfo `json:"details"`
}

func (s *Song) Score() float32 {
	var score float32
	defer func() {
		if err := recover(); err != nil {
			score = -1.0
		}
	}()
	score = float32(((s.Details.Likes - s.Details.Dislikes) / s.Details.Likes))
	return score
}