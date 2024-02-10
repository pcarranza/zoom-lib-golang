package zoom // Use this file for /recording endpoints

import (
	"fmt"
	"io"
	"net/http"
)

const (
	// RecordingTypeSharedScreenWithSpeakerViewCC is a shared screen with spearker view (CC) recording
	RecordingTypeSharedScreenWithSpeakerViewCC RecordingType = "shared_screen_with_speaker_view(CC)"
	// RecordingTypeSharedScreenWithSpeakerView is a shared screen with spearker view recording
	RecordingTypeSharedScreenWithSpeakerView RecordingType = "shared_screen_with_speaker_view"
	// RecordingTypeSharedScreenWithGalleryView is a shared screen with gallery view recording
	RecordingTypeSharedScreenWithGalleryView RecordingType = "shared_screen_with_gallery_view"
	// RecordingTypeSpeakerView is a speaker view recording
	RecordingTypeSpeakerView RecordingType = "speaker_view"
	// RecordingTypeGalleryView is a gallery view recording
	RecordingTypeGalleryView RecordingType = "gallery_view"
	// RecordingTypeSharedScreen is a shared screen recording
	RecordingTypeSharedScreen RecordingType = "shared_screen"
	// RecordingTypeAudioOnly is an audio only recording
	RecordingTypeAudioOnly RecordingType = "audio_only"
	// RecordingTypeAudioTranscript is an audio transcript recording
	RecordingTypeAudioTranscript RecordingType = "audio_transcript"
	// RecordingTypeChatFile is a chat file recording
	RecordingTypeChatFile RecordingType = "chat_file"
	// RecordingTypeTIMELINE is a timeline recording
	RecordingTypeTIMELINE RecordingType = "TIMELINE"
)

type (
	// RecordingType is the recording file type
	RecordingType string

	// RecordingFile represents a recordings file object
	RecordingFile struct {
		ID             string `json:"id"`
		MeetingID      string `json:"meeting_id"`
		RecordingStart *Time  `json:"recording_start"`
		RecordingEnd   *Time  `json:"recording_end"`
		FileType       string `json:"file_type"`
		FileSize       int    `json:"file_size"`
		FileExtension  string `json:"file_extension"`
		PlayURL        string `json:"play_url"`
		// The URL using which the recording file can be downloaded. To access a private or
		// password protected cloud recording, you must use a [Zoom JWT App Type]
		DownloadURL   string        `json:"download_url"`
		Status        string        `json:"status"`
		DeletedTime   *Time         `json:"deleted_time"`
		RecordingType RecordingType `json:"recording_type"`
	}

	// CloudRecordingMeeting represents a zoom meeting object
	CloudRecordingMeeting struct {
		UUID           string          `json:"uuid"`
		ID             int             `json:"id"`
		AccountID      string          `json:"account_id"`
		HostID         string          `json:"host_id"`
		Topic          string          `json:"topic"`
		StartTime      *Time           `json:"start_time"`
		Duration       int             `json:"duration"`
		TotalSize      int             `json:"total_size"`
		RecordingCount int             `json:"recording_count"`
		RecordingFiles []RecordingFile `json:"recording_files"`
	}
)

func (c *Client) DownloadRecordingFile(rf RecordingFile, w io.Writer) error {
	req, err := http.NewRequest("GET", rf.DownloadURL, nil)
	if err != nil {
		return err
	}

	req, err = c.addRequestAuth(req, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if _, err := io.Copy(w, resp.Body); err != nil {
			return fmt.Errorf("could not copy response body when downloading recording file: %w", err)
		}

		if err := resp.Body.Close(); err != nil {
			return fmt.Errorf("could not close body after copying: %w", err)
		}

	default:
		return fmt.Errorf("could not download recording file %s: %s", rf.ID, resp.Status)
	}

	return nil
}

// DeleteRecording calls DELETE /users/{user_id}/recordings/{recording_id}
// and deletes a cloud recording for the given user
func (c *Client) DeleteMeetingRecordings(opts DeleteMeetingRecordingsOptions) error {
	return c.requestV2(requestV2Opts{
		Method:        Delete,
		Path:          fmt.Sprintf(DeleteMeetingRecordingsPath, opts.MeetingID),
		URLParameters: &opts,
		HeadResponse:  true,
	})
}
