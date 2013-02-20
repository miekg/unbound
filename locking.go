package unbound

/*
#cgo LDFLAGS: -lunbound -lpthread
#include <openssl/ssl.h>
#include <openssl/err.h>
#include <unistd.h>

// Copied from Network Security with OpenSSL
// Copied from https://github.com/vadims/sslconn

#if defined(WIN32)
#define MUTEX_TYPE HANDLE
#define MUTEX_SETUP(x) (x) = CreateMutex(NULL, FALSE, NULL)
#define MUTEX_CLEANUP(x) CloseHandle(x)
#define MUTEX_LOCK(x) WaitForSingleObject((x), INFINITE)
#define MUTEX_UNLOCK(x) ReleaseMutex(x)
#define THREAD_ID GetCurrentThreadId( )
#elif defined(_POSIX_THREADS)
#include <pthread.h>
#define MUTEX_TYPE pthread_mutex_t
#define MUTEX_SETUP(x) pthread_mutex_init(&(x), NULL)
#define MUTEX_CLEANUP(x) pthread_mutex_destroy(&(x))
#define MUTEX_LOCK(x) pthread_mutex_lock(&(x))
#define MUTEX_UNLOCK(x) pthread_mutex_unlock(&(x))
#define THREAD_ID pthread_self( )
#else
#error You must define mutex operations appropriate for your platform!
#endif
static MUTEX_TYPE *mutex_buf = NULL;

static void locking_function(int mode, int n, const char *file, int line)
{
	if (mode & CRYPTO_LOCK)
		MUTEX_LOCK(mutex_buf[n]);
	else
		MUTEX_UNLOCK(mutex_buf[n]);
}

static unsigned long id_function(void)
{
	return ((unsigned long)THREAD_ID);
}

struct CRYPTO_dynlock_value {
	MUTEX_TYPE mutex;
};

static struct CRYPTO_dynlock_value *dyn_create_function(const char *file,
							int line)
{
	struct CRYPTO_dynlock_value *value;
	value =
	    (struct CRYPTO_dynlock_value *)
	    malloc(sizeof(struct CRYPTO_dynlock_value));
	if (!value)
		return NULL;
	MUTEX_SETUP(value->mutex);
	return value;
}

static void dyn_lock_function(int mode, struct CRYPTO_dynlock_value *l,
			      const char *file, int line)
{
	if (mode & CRYPTO_LOCK)
		MUTEX_LOCK(l->mutex);
	else
		MUTEX_UNLOCK(l->mutex);
}

static void dyn_destroy_function(struct CRYPTO_dynlock_value *l,
				 const char *file, int line)
{
	MUTEX_CLEANUP(l->mutex);
	free(l);
}

int ub_locking_init(void)
{
	int i;
	mutex_buf =
	    (MUTEX_TYPE *) malloc(CRYPTO_num_locks() * sizeof(MUTEX_TYPE));
	if (!mutex_buf)
		return 0;
	for (i = 0; i < CRYPTO_num_locks(); i++)
		MUTEX_SETUP(mutex_buf[i]);
	CRYPTO_set_id_callback(id_function);
	CRYPTO_set_locking_callback(locking_function);
	CRYPTO_set_dynlock_create_callback(dyn_create_function);
	CRYPTO_set_dynlock_lock_callback(dyn_lock_function);
	CRYPTO_set_dynlock_destroy_callback(dyn_destroy_function);
	return 1;
}

int ub_locking_cleanup(void)
{
	int i;
	if (!mutex_buf)
		return 0;
	CRYPTO_set_id_callback(NULL);
	CRYPTO_set_locking_callback(NULL);
	CRYPTO_set_dynlock_create_callback(NULL);
	CRYPTO_set_dynlock_lock_callback(NULL);
	CRYPTO_set_dynlock_destroy_callback(NULL);
	for (i = 0; i < CRYPTO_num_locks(); i++)
		MUTEX_CLEANUP(mutex_buf[i]);
	free(mutex_buf);
	mutex_buf = NULL;
	return 1;
}
*/
import "C"

func init() {
	C.ub_locking_init()
}
