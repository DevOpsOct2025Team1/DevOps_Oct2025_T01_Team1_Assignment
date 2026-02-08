import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { uploadFileInChunks } from './chunkedUpload';

describe('uploadFileInChunks', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn());
    localStorage.setItem('token', 'test_token');
  });

  afterEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    vi.unstubAllGlobals();
  });

  it('should upload file in chunks and track progress', async () => {
    const mockFile = new File(['x'.repeat(20 * 1024 * 1024)], 'test.mp4', {
      type: 'video/mp4',
    });

    const progressCallback = vi.fn();

    const fetchMock = vi.mocked(fetch);
    fetchMock
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          upload_id: 'test_upload_id',
          chunk_size: 10 * 1024 * 1024,
          total_parts: 2,
        }),
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({ etag: '"etag1"', part_number: 1 }),
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({ etag: '"etag2"', part_number: 2 }),
      } as Response)
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          file: { id: 'file_123', filename: 'test.mp4' },
        }),
      } as Response);

    const result = await uploadFileInChunks(mockFile, progressCallback);

    expect(result.file.id).toBe('file_123');
    expect(progressCallback).toHaveBeenCalled();
    expect(progressCallback).toHaveBeenCalledWith(50);
    expect(progressCallback).toHaveBeenCalledWith(100);
  });

  it('should abort upload on error', async () => {
    const mockFile = new File(['x'.repeat(20 * 1024 * 1024)], 'test.mp4');
    const progressCallback = vi.fn();

    const fetchMock = vi.mocked(fetch);
    fetchMock
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          upload_id: 'test_upload_id',
          chunk_size: 10 * 1024 * 1024,
          total_parts: 2,
        }),
      } as Response)
      .mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({ error: 'Server error' }),
      } as Response)
      .mockResolvedValueOnce({ ok: true } as Response);

    await expect(uploadFileInChunks(mockFile, progressCallback)).rejects.toThrow();

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/multipart/test_upload_id'),
      expect.objectContaining({ method: 'DELETE' })
    );
  });
});
