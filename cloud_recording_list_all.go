package zoom

import "fmt"

const (
	// ListAllRecordingsPath - v2 lists all recordings
	ListAllRecordingsPath = "/users/%s/recordings"

	DeleteMeetingRecordingsPath = "/meetings/%d/recordings"

	// TrashTypeMeetingRecordings list all meeting recordings from the trash. Default.
	TrashTypeMeetingRecordings TrashType = "meeting_recordings"
	// TrashTypeRecordingFile list all individual recording files from the trash
	TrashTypeRecordingFile TrashType = "recording_file"
)

type (
	// TrashType is the type of Cloud recording that you would like to retrieve from the trash
	TrashType string

	// ListAllRecordingsOptions contains options for ListAllRecordings.
	// NOTE: The query URL parser doesn't like non time.Time fields. It just
	// ignores this field if it's a zoom.Date or a zoom.Time. Instead use a
	// `string` for `From` and `To` - see below.
	ListAllRecordingsOptions struct {
		UserID        string `url:"-"`
		PageSize      *int   `url:"page_size,omitempty"`
		NextPageToken string `url:"next_page_token,omitempty"`
		Mc            string `url:"mc,omitempty"`
		Trash         bool   `url:"trash,omitempty"`
		// From is a YYYY-MM-DD string representing a date
		From string `url:"from,omitempty"`
		// To is a YYYY-MM-DD string representing a date
		To        string    `url:"to,omitempty"`
		TrashType TrashType `url:"trash_type,omitempty"`
	}

	// ListAllRecordingsResponse contains the response from a call to ListAllRecordings
	ListAllRecordingsResponse struct {
		From          *Date                   `json:"from"`
		To            *Date                   `json:"to"`
		PageCount     int                     `json:"page_count"`
		PageSize      int                     `json:"page_size"`
		TotalRecords  int                     `json:"total_records"`
		NextPageToken string                  `json:"next_page_token"`
		Meetings      []CloudRecordingMeeting `json:"meetings"`
	}

	DeleteMeetingRecordingsOptions struct {
		MeetingID int  `url:"-"`
		Delete    bool `url:",omitempty"` // By default it is trash, thus, on true, directly delete
	}
)

// ListAllRecordings calls /users/{user_id}/recordings endpoint
// and gets all cloud recordings for a user, using the default
// client.
func ListAllRecordings(opts ListAllRecordingsOptions) (ListAllRecordingsResponse, error) {
	return defaultClient.ListAllRecordings(opts)
}

// ListAllRecordings calls /users/{user_id}/recordings endpoint
// and gets all cloud recordings for a user, using the c client
//
// BEWARE: if you request more than 30 days, zoom will limit the response to
// that silently and automatically.
func (c *Client) ListAllRecordings(opts ListAllRecordingsOptions) (ListAllRecordingsResponse, error) {
	var ret = ListAllRecordingsResponse{}
	return ret, c.requestV2(requestV2Opts{
		Method:        Get,
		Path:          fmt.Sprintf(ListAllRecordingsPath, opts.UserID),
		URLParameters: &opts,
		Ret:           &ret,
	})
}
