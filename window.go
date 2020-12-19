package sliding

import (
	"time"

	"github.com/pkg/errors"
)

const (
	requestLimit = 100
)

var (
	ErrRequestLimit = errors.New("request limit is too high")
)

type Window struct {
	size         int
	full         bool
	begin        int
	nextWrite    int
	timeStamps   []int64
	currentWrite int
}

func NewWindow(size int) (*Window, error) {
	if size > requestLimit {
		return nil, errors.Wrapf(ErrRequestLimit, "max request limit: %d", requestLimit)
	}
	return &Window{
		size:       size,
		timeStamps: make([]int64, size),
	}, nil
}

// резервируем слот, возвращает false если слотов нет
func (w *Window) Reserve() bool {
	w.clean()
	return w.circularPut()
}

func (w *Window) IsAvailable() bool {
	w.clean()
	return w.checkAvailable()
}

func (w *Window) checkAvailable() bool {
	if w.begin == w.nextWrite && w.full {
		return false
	}
	return true
}

// простая реализация circular buffer
func (w *Window) circularPut() bool {
	if !w.checkAvailable() {
		return false
	}
	// записываем TS в конец
	w.currentWrite = w.nextWrite
	w.timeStamps[w.currentWrite] = getCurrentMillis()
	// увеличиваем указатель с учетом цикличности буфера
	w.nextWrite = (w.nextWrite + 1) % w.size
	// если указатель на следующую запись совпадает с началом буфера, значит он переполнен и дальше запись невозможна
	if w.begin == w.nextWrite {
		w.full = true
	}
	return true
}

// подчищает слоты, которые за пределами окна
func (w *Window) clean() {
	// край окна слева к текущему времени
	currentMillis := getCurrentMillis() - windowDuration()
	for {
		// очистили целиком
		if w.begin == w.currentWrite {
			return
		}
		// слева TS больше, чем край окна, значит больше ничего не учистить
		if w.timeStamps[w.begin] > currentMillis {
			return
		}
		// очищаем окно сдвигая курсор начала
		w.begin = (w.begin + 1) % w.size
	}
}

var getCurrentMillis = func() int64 {
	return time.Now().UnixNano() / 1000000
}

var windowDuration = func() int64 {
	return time.Second.Milliseconds()

}
