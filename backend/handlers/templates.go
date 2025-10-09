package handlers

import (
    "fmt"
    "os"
    "path/filepath"

    "farm-management-backend/database"
    "farm-management-backend/models"

    "github.com/gofiber/fiber/v2"
)

// resolvePath tries multiple locations to find an existing file path.
func resolvePath(p string) (string, error) {
    // If absolute and exists, use it
    if filepath.IsAbs(p) {
        if _, err := os.Stat(p); err == nil {
            return p, nil
        }
    } else {
        // Try repo-root style: backend/<relative>
        candidate := filepath.Join("backend", p)
        if _, err := os.Stat(candidate); err == nil {
            return candidate, nil
        }
        // Try as-is relative to current working dir
        if _, err := os.Stat(p); err == nil {
            return p, nil
        }
        // For common case of templates/<file>
        candidate2 := filepath.Join("templates", filepath.Base(p))
        if _, err := os.Stat(candidate2); err == nil {
            return candidate2, nil
        }
    }
    return "", fmt.Errorf("file not found: %s", p)
}

// DownloadTemplate serves Excel templates stored in backend/templates directory.
// Supported keys:
//  - upah        -> template_upah_24101780.xlsx
//  - tk          -> template_tk_24101780.xlsx
//  - tk_na       -> template_tk_na_24101780.xlsx
//  - koreksi_tk  -> template_koreksi_tk_24101780.xlsx
func DownloadTemplate() fiber.Handler {
    return func(c *fiber.Ctx) error {
        key := c.Params("key")

        // Try DB lookup if DB is available in context
        if dbLoc := c.Locals("db"); dbLoc != nil {
            if db, ok := dbLoc.(*database.DB); ok && db != nil {
                var tmpl models.Template
                if err := db.Where("key = ? AND is_active = ?", key, true).First(&tmpl).Error; err == nil {
                    // Build absolute path; allow stored absolute or relative paths
                    // Resolve actual path across different working dirs
                    path, perr := resolvePath(tmpl.FilePath)
                    if perr != nil {
                        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                            "error": perr.Error(),
                        })
                    }

                    // Set headers so browser downloads with the recorded filename and mimetype
                    if tmpl.MimeType != "" {
                        c.Set("Content-Type", tmpl.MimeType)
                    } else {
                        c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
                    }
                    c.Attachment(tmpl.FileName)

                    if err := c.SendFile(path, true); err != nil {
                        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                            "error": fmt.Sprintf("failed to send file: %v", err),
                        })
                    }
                    return nil
                }
            }
        }

        // Fallback: hardcoded mapping to files in backend/templates
        var filename string
        switch key {
        case "upah":
            filename = "template_upah_24101780.xlsx"
        case "tk":
            filename = "template_tk_24101780.xlsx"
        case "tk_na":
            filename = "template_tk_na_24101780.xlsx"
        case "koreksi_tk":
            filename = "template_koreksi_tk_24101780.xlsx"
        default:
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "template not found",
            })
        }

        // Try multiple candidate paths
        path, perr := resolvePath(filepath.Join("templates", filename))
        if perr != nil {
            // last resort: backend/templates/<file>
            path = filepath.Join("backend", "templates", filename)
        }
        c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
        c.Attachment(filename)
        if err := c.SendFile(path, true); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": fmt.Sprintf("failed to send file: %v", err),
            })
        }
        return nil
    }
}


