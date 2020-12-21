package sliding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCircular(t *testing.T) {
	size := 10
	w, err := NewWindow(size)
	require.NoError(t, err)
	for i := 0; i < 10; i++ {
		require.True(t, w.circularPut())
	}
	require.False(t, w.circularPut())
	require.False(t, w.circularPut())
}

func TestWindow_Put(t *testing.T) {
	var millis int64 = 5
	getCurrentMillis = func() int64 {
		return millis
	}
	var duration int64 = 2 // окно размер в 2 мс для простоты тестов
	windowDuration = func() int64 {
		return duration
	}
	size := 4
	w, err := NewWindow(size)
	require.NoError(t, err)
	// 4 запроса сразу разрешаем
	require.True(t, putN(w, 4))
	// пятый уже фейлится
	require.False(t, putN(w, 1))
	// на следующей мс можем снова 4, поэтому закидываем 2 в этом слоте
	millis = 6
	require.False(t, putN(w, 2))
	// и два в этом
	millis = 7
	require.True(t, putN(w, 2))
	// на этом времени освободилось только два слота, которые закидывали ранее
	millis = 8
	require.True(t, putN(w, 2))
	// третий уже не можем докинуть
	require.False(t, putN(w, 1))
	// пролистнули окно целиком, должно все очиститься
	millis += duration
	require.True(t, putN(w, 4))
}

func putN(w *Window, n int) bool {
	for i := 0; i < n; i++ {
		if !w.Reserve() {
			return false
		}
	}
	return true
}
