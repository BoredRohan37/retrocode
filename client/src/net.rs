//! Wire framing: 4-byte little-endian length prefix + UTF-8 JSON payload.
//! Bit 31 of the length word marks a zstd-compressed payload (parsed in M0,
//! compression enabled in M3). Max frame: 16 MiB.

pub const COMPRESSED_BIT: u32 = 0x8000_0000;
pub const MAX_FRAME: u32 = 16 * 1024 * 1024;

pub fn encode(payload: &[u8], compressed: bool) -> Vec<u8> {
    let mut len = payload.len() as u32;
    assert!(len <= MAX_FRAME, "frame too large");
    if compressed {
        len |= COMPRESSED_BIT;
    }
    let mut out = Vec::with_capacity(4 + payload.len());
    out.extend_from_slice(&len.to_le_bytes());
    out.extend_from_slice(payload);
    out
}

pub fn decode_header(h: &[u8; 4]) -> (u32, bool) {
    let raw = u32::from_le_bytes(*h);
    (raw & !COMPRESSED_BIT, raw & COMPRESSED_BIT != 0)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn roundtrip_plain() {
        let frame = encode(b"{\"m\":\"session/ping\"}", false);
        let (len, compressed) = decode_header(&frame[0..4].try_into().unwrap());
        assert_eq!(len as usize, frame.len() - 4);
        assert!(!compressed);
        assert_eq!(&frame[4..], b"{\"m\":\"session/ping\"}");
    }

    #[test]
    fn compressed_flag_survives() {
        let frame = encode(b"abc", true);
        let (len, compressed) = decode_header(&frame[0..4].try_into().unwrap());
        assert_eq!(len, 3);
        assert!(compressed);
    }

    #[test]
    #[should_panic]
    fn rejects_oversized_frame() {
        let big = vec![0u8; (MAX_FRAME + 1) as usize];
        encode(&big, false);
    }
}
