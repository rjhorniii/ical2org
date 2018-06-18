;;; ical2org-integration.el --- Provide integration for ical2org and various mailers

;;; Commentary:
;; 


;;; Code:

(defvar ical2org-org-file "~/org/ical-import.org")

(defconst +ical2org-executable+
  (executable-find "ical2org"))

(defvar ical2org-gnus-schedule-key "i")
(defvar ical2org-gnus-view-key "C-M-i")

(defun ical2org-mu4e-schedule (msg attachnum)
  "Schedule the ATTACHNUM in MSG."
  (mu4e-view-pipe-attachment msg attachnum
                             (format "%s -count -d=%s -a=%s -"
                                     +ical2org-executable+
                                     ical2org-org-file
                                     ical2org-org-file)))

(defun ical2org-mu4e-view-as-org (msg attachnum)
  "View the ATTACHNUM in MSG as Org."
  (mu4e-view-pipe-attachment msg attachnum
                             (format "%s -"
                                     (executable-find "ical2org"))))

(defun ical2org-mu4e-insinuate ()
  "Enable mu4e ical2org integration"
  (eval-after-load 'mu4e
    `(progn
       (add-to-list 'mu4e-view-attachment-actions '("schedule appointment" . ical2org-mu4e-schedule) t)
       (add-to-list 'mu4e-view-attachment-actions '("view appointment" . ical2org-mu4e-view-as-org) t))))

(defun ical2org--gnus-fetch-vcal ()
  "Fetch the vcal data from the current message"
  (save-excursion
    (let (invitation
          beginning
          ending
          (content-buffer (current-buffer)))
      (goto-char (point-min))
      (if (not (re-search-forward "^BEGIN:VCALENDAR" (point-max) t))
          (error "Unable to find start of VCalendar")
        (setq beginning (match-beginning 0)))
      (setf ending (re-search-forward "^END:VCALENDAR\n?" (point-max) t))
      (when (not ending)
        (error "Unable to find end of VCalendar"))
      (buffer-substring beginning ending))))

(defun ical2org-gnus-schedule ()
  "Schedule the vcal data in the current message."
  (interactive)
  (save-excursion
    (when (equal major-mode 'gnus-article-mode)
      (gnus-article-show-summary))
    (when (equal major-mode 'gnus-summary-mode)
      (gnus-summary-show-article)
      (gnus-summary-select-article-buffer)
      (gnus-mime-view-all-parts)
      (let ((ical-data (ical2org--gnus-fetch-vcal)))
        (with-temp-buffer
          (insert ical-data)
          (call-process-region (point-min) (point-max)
                               +ical2org-executable+
                               nil nil nil
                               "-count"
                               (format "-d=%s" ical2org-org-file)
                               (format "-a=%s" ical2org-org-file)
                               "-")))
      (gnus-summary-show-article))))

(defun ical2org-gnus-view ()
  "View the vcal data in the current message as Org."
  (interactive)
  (save-excursion
    (when (equal major-mode 'gnus-article-mode)
      (gnus-article-show-summary))
    (when (equal major-mode 'gnus-summary-mode)
      (gnus-summary-show-article) ; let Gnus decode b64, QP
      (gnus-summary-select-article-buffer)
      (gnus-mime-view-all-parts) ; looking for text/calendar
      (let ((ical-data (ical2org--gnus-fetch-vcal)))
        (with-temp-buffer
          (insert ical-data)
          (call-process-region (point-min) (point-max)
                               +ical2org-executable+
                               t t nil
                               "-")
          (let ((string (buffer-string)))
            (with-output-to-temp-buffer "*ical2org-view*"
              (princ string)))))
      (gnus-summary-show-article))))

(defun ical2org-gnus-insinuate ()
  "Enable ical2org GNUS integration"
  (eval-after-load 'gnus-summary
    `(progn
       (define-key gnus-summary-mode-map (kbd ,ical2org-gnus-schedule-key) #'ical2org-gnus-schedule)
       (define-key gnus-summary-mode-map (kbd ,ical2org-gnus-view-key) #'ical2org-gnus-view)))
  (eval-after-load 'gnus
    (define-key gnus-article-mode-map (kbd ,ical2org-gnus-schedule-key) #'ical2org-gnus-schedule)
    (define-key gnus-article-mode-map (kbd ,ical2org-gnus-view-key) #'ical2org-gnus-view)))

(provide 'ical2org-integration)

;;; ical2org-integration.el ends here
